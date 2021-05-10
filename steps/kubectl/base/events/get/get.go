package get

import (
	"encoding/json"
	"fmt"

	"github.com/stackpulse/steps-sdk-go/log"
	base2 "github.com/stackpulse/steps/kubectl/base"
)

const ObjectNameJsonKey = "objectName"

var ObjectNameFilterAlreadyExists = fmt.Errorf("can't add more than one filter for object name")

type Args struct {
	base2.Args
	ObjectName   string `env:"OBJECT_NAME"`
	ObjectType   string `env:"OBJECT_TYPE"`
	NameContains string `env:"NAME_CONTAINS"`
	NameExact    string `env:"NAME_EXACT"`
}

func (a *Args) addFilters() error {
	if a.NameContains != "" {
		if _, ok := a.FilterContainsParsed[ObjectNameJsonKey]; ok {
			return ObjectNameFilterAlreadyExists
		}
		a.FilterContainsParsed[ObjectNameJsonKey] = a.NameContains
	}

	if a.NameExact != "" {
		if _, ok := a.FilterEqualsParsed[ObjectNameJsonKey]; ok {
			return ObjectNameFilterAlreadyExists
		}
		a.FilterEqualsParsed[ObjectNameJsonKey] = a.NameExact
	}
	return nil
}

type GetEvents struct {
	Args *Args
	kctl *base2.KubectlStep
}

func NewGetEvents(args *Args) (*GetEvents, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}
	kctl, err := base2.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	if err := args.addFilters(); err != nil {
		return nil, err
	}

	return &GetEvents{
		Args: args,
		kctl: kctl,
	}, nil
}

func (e *GetEvents) Get() (output []byte, exitCode int, err error) {
	cmdArgs := []string{"get", "events"}

	if e.Args.ObjectName != "" {
		cmdArgs = append(cmdArgs, "--field-selector", fmt.Sprintf("involvedObject.name=%s", e.Args.ObjectName))
	}

	if e.Args.ObjectType != "" {
		cmdArgs = append(cmdArgs, "--field-selector", fmt.Sprintf("involvedObject.kind=%s", e.Args.ObjectType))
	}

	return e.kctl.Execute(cmdArgs)
}

func (e *GetEvents) Parse(output []byte) (string, error) {
	return e.kctl.ParseOutput(output, parsingConfiguration)
}

func (e *GetEvents) ParseObject(output []byte) (*Events, error) {
	parsed, err := e.Parse(output)
	if err != nil {
		return nil, err
	}

	ret := &Events{}
	if err := json.Unmarshal([]byte(parsed), ret); err != nil {
		log.Logln("JSON unmarshaling failed. Original output:\n%s", parsed)
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return ret, nil
}