package delete

import (
	"fmt"
	"strconv"
	"time"

	base2 "github.com/stackpulse/steps/kubectl/base"
)

type Args struct {
	base2.Args
	ResourceType   string        `env:"RESOURCE_TYPE,required"`
	ResourcesNames []string      `env:"RESOURCES_NAMES"`
	DeleteAll      bool          `env:"DELETE_ALL" envDefault:"false"`
	GracePeriod    int           `env:"GRACE_PERIOD" envDefault:"-1"`
	IgnoreNotFound bool          `env:"IGNORE_NOT_FOUND" envDefault:"false"`
	Cascade        bool          `env:"CASCADE" envDefault:"true"`
	Force          bool          `env:"FORCE" envDefault:"false"`
	Timeout        time.Duration `env:"TIMEOUT" envDefault:"0s"`
}

type Delete struct {
	Args *Args
	kctl *base2.KubectlStep
}

func validate(args *Args) error {
	if len(args.ResourcesNames) == 0 && !args.DeleteAll {
		return fmt.Errorf("must sprcify RESOURCES_NAMES or DELETE_ALL")
	}
	return nil
}

func NewDelete(args *Args) (*Delete, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}

	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if err := validate(args); err != nil {
		return nil, err
	}

	return &Delete{
		Args: args,
		kctl: kctl,
	}, nil
}

func (d *Delete) getReourceToDelete() []string {
	if d.Args.DeleteAll {
		return []string{"--all"}
	}

	toDelete := make([]string, 0, len(d.Args.ResourcesNames))
	for _, v := range d.Args.ResourcesNames {
		// Deleting empty strings
		if v == "" {
			continue
		}
		toDelete = append(toDelete, v)
	}
	return toDelete
}

func (d *Delete) Delete() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"delete", d.Args.ResourceType}
	cmdArgs = append(cmdArgs, d.getReourceToDelete()...)
	cmdArgs = append(cmdArgs, "--grace-period", strconv.Itoa(d.Args.GracePeriod))
	cmdArgs = append(cmdArgs, "--ignore-not-found", strconv.FormatBool(d.Args.IgnoreNotFound))
	cmdArgs = append(cmdArgs, "--cascade", strconv.FormatBool(d.Args.Cascade))
	cmdArgs = append(cmdArgs, "--timeout", d.Args.Timeout.String())

	if d.Args.Force {
		cmdArgs = append(cmdArgs, "--force")
	}

	// Ignoring format because delete has not json format
	return d.kctl.Execute(cmdArgs, base2.IgnoreFormat)
}
