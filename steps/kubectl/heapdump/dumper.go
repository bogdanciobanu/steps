package main

import (
	"fmt"
	"strings"

	"github.com/stackpulse/steps/kubectl/base"
)

type Dumper interface {
	Init() error
	Prerequisites() error
	Dump() (string, error)
	Cleanup() error
}

func NewDumper(runtime string, kubectl *base.KubectlStep, pod PodInfo) (Dumper, error) {
	switch strings.ToLower(runtime) {
	case "jvm":
		return NewJAttach(kubectl, pod), nil
	default:
	}

	return nil, fmt.Errorf("unsupported runtime: %s", runtime)
}
