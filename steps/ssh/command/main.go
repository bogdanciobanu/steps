package main

import (
	"encoding/json"
	"fmt"

	sdkExec "github.com/stackpulse/steps-sdk-go/exec"
	"github.com/stackpulse/steps-sdk-go/step"
	sshBase "github.com/stackpulse/steps/ssh/base"
)

type Outputs struct {
	CommandOutput string `json:"output"`
	step.Outputs
}

type SSHCommand struct {
	args sshBase.SSHArgs
}

func (s *SSHCommand) Init() error {
	var err error
	s.args, err = sshBase.ParseArgs()
	return err
}

func (s *SSHCommand) marshalOutput(output string) []byte {
	outputJSON, _ := json.Marshal(Outputs{
		CommandOutput: output,
	})
	return outputJSON
}

func (s *SSHCommand) Run() (int, []byte, error) {
	var sshCmd string
	var sshArgs []string

	sshCmd, sshArgs, err := s.args.BuildCommand()
	if err != nil {
		return step.ExitCodeFailure, s.marshalOutput(""), fmt.Errorf("buildCommand: %w", err)
	}

	output, exitCode, err := sdkExec.Execute(sshCmd, sshArgs)
	sshBase.PrintLog()
	sshBase.PrintOutput(output)
	marshaledOuptut := s.marshalOutput(string(output))
	if err != nil {
		return exitCode, marshaledOuptut, fmt.Errorf("execute ssh: %w", err)
	}

	return step.ExitCodeOK, marshaledOuptut, nil
}

func main() {
	step.Run(&SSHCommand{})
}
