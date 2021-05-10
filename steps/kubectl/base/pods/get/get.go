package get

import (
	"encoding/json"
	"fmt"
	"github.com/stackpulse/steps-sdk-go/log"

	"github.com/stackpulse/steps-sdk-go/filter"
	"github.com/stackpulse/steps/kubectl/base"
	"maze.io/x/duration.v1"
)

type Args struct {
	base.Args
	PodNames     []string          `env:"POD_NAMES"`
	NameContains string            `env:"NAME_CONTAINS"`
	NameExact    string            `env:"NAME_EXACT"`
	SineStr      string            `env:"SINCE"`
	Since        duration.Duration `env:"-"`
}

var PodNameFilterAlreadyExists = fmt.Errorf("can't add more than one filter for pod name")

func (a *Args) addFilters() error {
	if a.NameContains != "" {
		if _, ok := a.FilterContainsParsed[PodNameJsonKey]; ok {
			return PodNameFilterAlreadyExists
		}
		a.FilterContainsParsed[PodNameJsonKey] = a.NameContains
	}

	if a.NameExact != "" {
		if _, ok := a.FilterEqualsParsed[PodNameJsonKey]; ok {
			return PodNameFilterAlreadyExists
		}
		a.FilterEqualsParsed[PodNameJsonKey] = a.NameExact
	}
	return nil
}

type GetPods struct {
	Args *Args
	kctl *base.KubectlStep
}

func NewGetPods(args *Args) (*GetPods, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if args.SineStr != "" {
		args.Since, err = duration.ParseDuration(args.SineStr)
		if err != nil {
			return nil, fmt.Errorf("can't parse since as duration: %w", err)
		}
	}

	if err := args.addFilters(); err != nil {
		return nil, err
	}

	return &GetPods{
		Args: args,
		kctl: kctl,
	}, nil
}

func (n *GetPods) Get() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"get", "pods"}
	cmdArgs = append(cmdArgs, n.Args.PodNames...)
	return n.kctl.Execute(cmdArgs)
}

func (n *GetPods) Parse(output []byte) (string, error) {
	var extraFilters []filter.JSONFilter
	if n.Args.Since > 0 {
		extraFilters = append(extraFilters, filter.JSONFilter{
			Matching:       map[string]string{"age": n.Args.Since.String()},
			ComparisonFunc: filter.DurationSmallerEqual,
		})
	}

	return n.kctl.ParseOutput(output, parsingConfiguration, extraFilters...)
}

func (n *GetPods) ParseObject(output []byte) (*Pods, error) {
	parsed, err := n.Parse(output)
	if err != nil {
		return nil, err
	}

	ret := &Pods{}
	if err := json.Unmarshal([]byte(parsed), ret); err != nil {
		log.Logln("JSON unmarshaling failed. Original output:\n%s", parsed)
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return ret, nil
}

