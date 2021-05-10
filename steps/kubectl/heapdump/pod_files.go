package main

import (
	"fmt"
	"github.com/google/shlex"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/kubectl/base"
)

type podFiles struct {
	localPath  string
	remotePath string
	PodInfo
}

func (p podFiles) command(isDownload bool) string {
	result := "cp "

	podPath := fmt.Sprintf("%s:%s", p.PodName, p.remotePath)
	if p.Namespace != "" {
		podPath = fmt.Sprintf("%s/%s", p.Namespace, podPath)
	}

	if isDownload {
		result += fmt.Sprintf("%s %s", podPath, p.localPath)
	} else {
		result += fmt.Sprintf("%s %s", p.localPath, podPath)
	}

	if p.Container != "" {
		result += " -c " + p.Container
	}

	return result
}

func (p podFiles) copy(kbl *base.KubectlStep, isDownload bool) ([]byte, int, error) {
	copyCmd := p.command(isDownload)
	cmdArgs, err := shlex.Split(copyCmd)
	if err != nil {
		return nil, step.ExitCodeFailure, fmt.Errorf("split copy command: %w", err)
	}

	return kbl.Execute(cmdArgs, base.IgnoreFormat)
}

func (p podFiles) upload(kbl *base.KubectlStep) ([]byte, int, error) {
	return p.copy(kbl, false)
}

func (p podFiles) download(kbl *base.KubectlStep) ([]byte, int, error) {
	return p.copy(kbl, true)
}
