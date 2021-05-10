package get

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	base2 "github.com/stackpulse/steps/kubectl/base"
)

type Args struct {
	base2.Args
	ApplyFile    string `env:"APPLY_FILE_PATH"`
	ApplyContent string `env:"APPLY_CONTENT"`
}

type ApplyFile struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewApply(args *Args) (*ApplyFile, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	return &ApplyFile{
		Args: args,
		kctl: kctl,
	}, nil
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

func (n *ApplyFile) Run() (output []byte, exitCode int, err error) {
	if n.Args.ApplyContent != "" {
		err = decodeAndWrite(n.Args.ApplyContent, "/tmp/data")
		if err != nil {
			fmt.Printf("errr! %s\n", err.Error())
		}
		n.Args.ApplyFile = "/tmp/data"
	}

	cmdArgs := []string{"apply", "-f"}
	cmdArgs = append(cmdArgs, n.Args.ApplyFile)

	return n.kctl.Execute(cmdArgs)
}

func (n *ApplyFile) Apply() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"apply", "-f"}
	cmdArgs = append(cmdArgs, n.Args.ApplyFile)
	return n.kctl.Execute(cmdArgs)
}
