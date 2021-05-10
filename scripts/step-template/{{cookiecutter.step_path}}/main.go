package main

import (
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type {{cookiecutter.normalized_step_name}} struct {
    {{cookiecutter.step_args}}
}

type output struct {
    {{cookiecutter.step_output}}
	step.Outputs
}

func (s *{{cookiecutter.normalized_step_name}}) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *{{cookiecutter.normalized_step_name}}) Run() (exitCode int, output []byte, err error) {
    // Replace This Comment with Step Logic

	return step.ExitCodeOK, []byte{}, nil
}

func main() {
	step.Run(&{{cookiecutter.normalized_step_name}}{})
}
