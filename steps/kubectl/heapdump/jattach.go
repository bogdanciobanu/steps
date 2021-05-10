package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps/kubectl/base"
)

type jattach struct {
	LocalPath            string `env:"JATTACH_LOCAL_PATH" envDefault:"/jattach"`
	RemotePath           string `env:"JATTACH_REMOTE_PATH" envDefault:"/tmp/jattach"`
	OutputDumpPath       string `env:"JATTACH_OUT_DUMP_PATH" envDefault:"/tmp/heapdump.hprof"`
	DumpScriptRemotePath string `env:"JATTACH_DUMP_SCRIPT_PATH" envDefault:"/tmp/jattach-dump.sh"`
	PulledOutputDumpPath string `env:"JATTACH_PULLED_DUMP_PATH" envDefault:"/tmp/heapdump.hprof"`
	kubectl              *base.KubectlStep
	pod                  PodInfo
}

const (
	dumpScriptLocalPath = "/tmp/jattach-dump.sh"
)

func NewJAttach(kubectl *base.KubectlStep, pod PodInfo) *jattach {
	return &jattach{
		kubectl: kubectl,
		pod:     pod,
	}
}

func (j *jattach) Init() error {
	return env.Parse(j)
}

func (j *jattach) Prerequisites() error {
	dumpScriptFile, err := os.Create(dumpScriptLocalPath)
	if err != nil {
		return fmt.Errorf("open dump script file at %q: %w", dumpScriptLocalPath, err)
	}
	defer closeWithLog(dumpScriptFile)

	err = j.writeDumpScript(dumpScriptFile)
	if err != nil {
		return fmt.Errorf("generate dump script: %w", err)
	}

	// upload required files
	for local, remote := range map[string]string{
		j.LocalPath:         j.RemotePath,
		dumpScriptLocalPath: j.DumpScriptRemotePath,
	} {
		out, _, err := podFiles{
			localPath:  local,
			remotePath: remote,
			PodInfo:    j.pod,
		}.upload(j.kubectl)
		if err != nil {
			return fmt.Errorf("upload file %q to %q due to %s: %w", local, remote, string(out), err)
		}
		j.kubectl.Debugln(string(out))
	}

	return nil
}

func (j *jattach) Dump() (string, error) {
	dumpCommand := fmt.Sprintf("chmod +x %q && %q", j.DumpScriptRemotePath, j.DumpScriptRemotePath)
	err := podExec(j.kubectl, j.pod, dumpCommand)
	if err != nil {
		return "", fmt.Errorf("execute dump command on pod: %w", err)
	}

	src, dst := j.OutputDumpPath, j.PulledOutputDumpPath
	out, _, err := podFiles{
		localPath:  src,
		remotePath: dst,
		PodInfo:    j.pod,
	}.download(j.kubectl)
	if err != nil {
		return "", fmt.Errorf("download file from '%s' to '%s' due to %s: %w", src, dst, string(out), err)
	}

	j.kubectl.Debugln(string(out))

	return dst, nil
}

func (j *jattach) Cleanup() error {
	files := []string{
		j.DumpScriptRemotePath,
		j.OutputDumpPath,
		j.RemotePath,
	}
	return podExec(j.kubectl, j.pod, fmt.Sprintf("rm -f %s", strings.Join(files, " ")))
}

func (j *jattach) writeDumpScript(writer io.Writer) error {
	commands := []string{
		fmt.Sprintf("chmod +x %q", j.RemotePath),
		fmt.Sprintf("%s $(pidof -s java) dumpheap %q", j.RemotePath, j.OutputDumpPath),
	}

	for _, c := range commands {
		_, err := writer.Write([]byte(c + "\n"))
		if err != nil {
			return fmt.Errorf("write command: '%s': %w", c, err)
		}
	}
	return nil
}
