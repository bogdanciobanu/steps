package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps-sdk-go/storage"
	"github.com/stackpulse/steps/kubectl/base"
)

type memoryDump struct {
	base.Args
	PodInfo
	Storage        storage.Storage
	Runtime        string `env:"RUNTIME,required"`

	kubectl        *base.KubectlStep
	dumper         Dumper
}

type PodInfo struct {
	PodName        string `env:"POD_NAME,required"`
	Container      string `env:"CONTAINER"`
	Namespace      string `env:"NAMESPACE"`
	ContainerShell string `env:"CONTAINER_SHELL" envDefault:"bash"`
}

type output struct {
	DumpFileURI string `json:"dump_file_uri"`
	step.Outputs
}

func closeWithLog(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Printf("Failed to perform close: %v", err)
	}
}

func podExec(kubectl *base.KubectlStep, pod PodInfo, command string) error {
	execCmd := []string{"exec", pod.PodName}
	if pod.Container != "" {
		execCmd = append(execCmd, "-c", pod.Container)
	}

	baseCmd, err := kubectl.BaseCommand(base.IgnoreFormat)
	if err != nil {
		return fmt.Errorf("kubectl base command: %w", err)
	}
	execCmd = append(append(execCmd, baseCmd...), "--", pod.ContainerShell, "-c", command)

	out, _, err := kubectl.Execute(execCmd, base.IgnoreFormat)
	if err != nil {
		return fmt.Errorf("execute command due to %s: %w", string(out), err)
	}
	kubectl.Debugln(string(out))
	return nil
}

func (m *memoryDump) Init() error {
	err := env.Parse(m)
	if err != nil {
		return err
	}

	kubectl, err := base.NewKubectlStep(&m.Args, true)
	if err != nil {
		return fmt.Errorf("parse kubectl step params: %w", err)
	}
	m.kubectl = kubectl

	dumper, err := NewDumper(m.Runtime, kubectl, m.PodInfo)
	if err != nil {
		return fmt.Errorf("create dumper: %w", err)
	}
	if err = dumper.Init(); err != nil {
		return fmt.Errorf("init dumper: %w", err)
	}
	m.dumper = dumper

	err = m.Storage.Init()
	if err != nil {
		return fmt.Errorf("init storage: %w", err)
	}

	return nil
}

func (m *memoryDump) Run() (int, []byte, error) {
	var err error
	defer func() {
		if cleanErr := m.dumper.Cleanup(); cleanErr != nil {
			log.Println("Cleanup failed due to: ", cleanErr)
		}
	}()
	err = m.dumper.Prerequisites()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("install prerequisites: %w", err)
	}

	dumpFile, err := m.dumper.Dump()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("dump: %w", err)
	}

	dumpFileKey := fmt.Sprintf("heapdump_%d.hprof", time.Now().Unix())
	err = m.Storage.UploadFiles(context.Background(), []storage.UploadFile{{Source: dumpFile, Key: dumpFileKey}})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("upload dump file: %w", err)
	}

	jsonOutput, err := json.Marshal(&output{
		DumpFileURI: m.Storage.Bucket + "/" + dumpFileKey,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	env.SetFormatter(env.JsonFormat, true)
	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&memoryDump{})
}
