package troubleshoot

import (
	"encoding/json"
	"fmt"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/kubectl/base"
	events "github.com/stackpulse/steps/kubectl/base/events/get"
)

type Args struct {
	base.Args
}

type Troubleshoot struct {
	Args *Args
	kctl *base.KubectlStep
}

type TroubleshootOutput map[string][]events.Event

func NewTroubleshoot(args *Args) (*Troubleshoot, error) {
	parse := false
	if args == nil {
		parse = true
		args = &Args{}
	}

	kctl, err := base.NewKubectlStep(args, parse)
	if err != nil {
		return nil, err
	}

	args.FieldSelector = append(args.FieldSelector)

	return &Troubleshoot{
		Args: args,
		kctl: kctl,
	}, nil
}

func (t *Troubleshoot) getEvents() (*events.Events, error) {
	t.Args.FieldSelector = append(t.Args.FieldSelector, "type!=Normal")
	args := &events.Args{
		Args: t.Args.Args,
	}
	eventsGet, err := events.NewGetEvents(args)
	if err != nil {
		return nil, fmt.Errorf("new get events: %w", err)
	}

	output, _, err := eventsGet.Get()
	if err != nil {
		return nil, fmt.Errorf("get events: %w. Output: %s", err, string(output))
	}

	return eventsGet.ParseObject(output)
}

func (t *Troubleshoot) getInitialResourcessMap() TroubleshootOutput {
	return map[string][]events.Event{
		"Pod":                     nil,
		"Deployment":              nil,
		"ReplicaSet":              nil,
		"StatefulSet":             nil,
		"HorizontalPodAutoscaler": nil,
		"CronJob":                 nil,
		"Job":                     nil,
		"DaemonSet":               nil,
		"Service":                 nil,
		"PersistentVolumeClaim":   nil,
		"Endpoints":               nil,
		"Ingress":                 nil,
	}
}

func (t *Troubleshoot) Run() ([]byte, int, error) {
	namespace := fmt.Sprintf("namespace %s", t.Args.Namespace)
	if t.Args.AllNamespaces {
		namespace = "all namespaces"
	}
	log.Logln("Getting troubleshooting information for %s", namespace)

	currentEvents, err := t.getEvents()
	if err != nil {
		return nil, step.ExitCodeFailure, err
	}

	resourceMap := t.getInitialResourcessMap()
	for _, event := range currentEvents.Items {
		resourceMap[event.ObjectType] = append(resourceMap[event.ObjectType], event)
	}

	output, err := json.Marshal(resourceMap)
	if err != nil {
		return output, step.ExitCodeFailure, nil
	}
	return output, step.ExitCodeOK, nil
}

func (t *Troubleshoot) ParseObject(output []byte) (TroubleshootOutput, error) {
	var troubleshootOutput TroubleshootOutput
	if err := json.Unmarshal(output, &troubleshootOutput); err != nil {
		return nil, fmt.Errorf("unmarshal output: %w", err)
	}

	return troubleshootOutput, nil
}
