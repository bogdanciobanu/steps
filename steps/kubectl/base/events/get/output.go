package get

import (
	base2 "github.com/stackpulse/steps/kubectl/base"
	"time"
)

var parsingConfiguration = map[string]*base2.JsonParseConfig{
	"objectType":      {ParseFunc: base2.JsonPathStringParser, Args: []string{"involvedObject.kind"}},
	"objectName":      {ParseFunc: base2.JsonPathStringParser, Args: []string{"involvedObject.name"}},
	"objectNamespace": {ParseFunc: base2.JsonPathStringParser, Args: []string{"involvedObject.namespace"}},
	"lastTimestamp":   {ParseFunc: base2.JsonPathStringParser, Args: []string{"lastTimestamp"}},
	"message":         {ParseFunc: base2.JsonPathStringParser, Args: []string{"message"}},
	"reason":          {ParseFunc: base2.JsonPathStringParser, Args: []string{"reason"}},
	"type":            {ParseFunc: base2.JsonPathStringParser, Args: []string{"type"}},
}

type Event struct {
	LastTimestamp   time.Time
	Message         string
	ObjectName      string
	ObjectNamespace string
	ObjectType      string
	Reason          string
	Type            string
}

type Events struct {
	Items []Event
}
