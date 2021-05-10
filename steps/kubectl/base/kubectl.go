package base

import (
	"os/exec"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
)

type KubectlStep struct {
	StepArgs *Args
}

func NewKubectlStep(args BaseArgs, parse bool) (*KubectlStep, error) {
	if parse {
		if err := Parse(args); err != nil {
			return nil, err
		}
	}
	return &KubectlStep{StepArgs: args.BaseArgs()}, nil
}

func inIgnoredArray(val IgnoredArgs, ignoredArgs []IgnoredArgs) bool {
	for _, i := range ignoredArgs {
		if val == i {
			return true
		}
	}
	return false
}

func (k *KubectlStep) Debugln(msg string, args ...interface{}) {
	if k.StepArgs.Debug {
		log.Logln(msg, args...)
	}
}

// Return base kubectl arguments (such as -n <namespace>, auth settings, etc..)
func (k *KubectlStep) BaseCommand(ignoreFields ...IgnoredArgs) ([]string, error) {
	authMethod, err := k.getAuthMethod()
	if err != nil {
		return nil, err
	}

	args := []string{authMethod}

	format := k.getFormat()
	if !inIgnoredArray(IgnoreFormat, ignoreFields) {
		args = append(args, format...)
	}
	if !inIgnoredArray(IgnoreNamespace, ignoreFields) {
		args = append(args, k.getNamespace()...)
	}

	if !inIgnoredArray(IgnoreFieldSelector, ignoreFields) {
		for _, fieldSelector := range k.StepArgs.FieldSelector {
			args = append(args, "--field-selector", fieldSelector)
		}
	}

	return args, nil
}

// Run kubectl with given arguments
func (k *KubectlStep) Run(args []string) (output []byte, exitCode int, err error) {
	cmd := exec.Command("kubectl", args...)
	k.Debugln("About to run kubectl with args: %+v", args)
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			return cmdOutput, exiterr.ExitCode(), err
		}
		return cmdOutput, 1, err
	}

	return cmdOutput, 0, nil
}

// Combine BaseCommand() and Run() to execute kubectl command
// The extraArgs are arguments to append to kubectl beside the base arguments
// The base arguments will come AFTER the extraArgs
func (k *KubectlStep) Execute(extraArgs []string, ignoreFields ...IgnoredArgs) (output []byte, exitCode int, err error) {
	baseCommands, err := k.BaseCommand(ignoreFields...)
	if err != nil {
		return nil, 1, err
	}

	kctl := append(extraArgs, baseCommands...)
	return k.Run(kctl)
}

func (k *KubectlStep) getNamespace() []string {
	if k.StepArgs.AllNamespaces {
		return []string{"--all-namespaces"}
	}
	return []string{"-n", k.StepArgs.Namespace}
}

func (k *KubectlStep) getFormat() []string {
	if k.StepArgs.Format != "json" {
		env.SetFormatter(env.RawFormat, false)
	}

	if k.StepArgs.Format == "default" {
		return []string{}
	}
	return []string{"-o", k.StepArgs.Format}
}
