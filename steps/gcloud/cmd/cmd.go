package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

const (
	gcloudAuthTmpFile = "/tmp/auth.json"
)

type Outputs struct {
	Message string `json:"message"`
}

type GcloudCommand struct {
	args Args
}

type Args struct {
	Command    string `env:"COMMAND"`
	GcloudAuth string `env:"AUTH_CODE"`
}

func decodeAndWrite(base64Encoded, destination string) error {
	decoded, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(destination, decoded, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func cmdToSlice(command string) []string {
	splitCommand := strings.Split(command, "--")
	splittedBaseCommand := strings.Split(splitCommand[0], " ")
	result := []string{}
	var temp string
	for _, s := range splittedBaseCommand {
		if s != "" {
			result = append(result, s)
		}
	}
	for _, s := range splitCommand[1:] {
		temp = fmt.Sprintf("--%s", strings.TrimSpace(s))
		result = append(result, temp)
	}
	return result
}

func buildOutput(message string) ([]byte, error) {
	outputJSON, err := json.Marshal(Outputs{
		Message: message,
	})
	if err != nil {
		return nil, fmt.Errorf("json marshal error: %w", err)
	}
	return outputJSON, nil
}

func (cmd *GcloudCommand) runGcloudAuth() error {
	if _, err := exec.LookPath("gcloud"); err != nil {
		return fmt.Errorf("gcloud command not found (are you running <step>-gke variation?)")
	}

	if err := decodeAndWrite(cmd.args.GcloudAuth, gcloudAuthTmpFile); err != nil {
		return fmt.Errorf("can't write decoded gcloud auth: %w", err)
	}

	if _, err := exec.Command("gcloud", "auth", "activate-service-account", "--key-file", gcloudAuthTmpFile).CombinedOutput(); err != nil {
		return fmt.Errorf("can't run gcloud activate-service-account command: %w", err)
	}

	return nil
}

func (cmd *GcloudCommand) Init() error {
	var err error
	err = env.Parse(&cmd.args)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	err = cmd.runGcloudAuth()
	if err != nil {
		return fmt.Errorf("gcloud auth error: %w", err)
	}
	return nil
}

func (cmd *GcloudCommand) Run() (int, []byte, error) {
	command := cmdToSlice(cmd.args.Command)
	execution := exec.Command("gcloud", command...)
	cmdOutput, err := execution.CombinedOutput()
	if err != nil {
		log.Debug(string(cmdOutput))
		output, _ := buildOutput(string(cmdOutput))
		return step.ExitCodeFailure, output, err
	}
	output, _ := buildOutput("command ran succesfully")
	return step.ExitCodeOK, output, nil
}

func main() {
	step.Run(&GcloudCommand{})
}
