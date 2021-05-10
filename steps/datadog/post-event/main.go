package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"

	datadog "github.com/DataDog/datadog-api-client-go/api/v1/datadog"
)

const stackpulseSource = "STACKPULSE"

type EventPost struct {
	ApiKey         string   `env:"DD_API_KEY,required"`
	Site           string   `env:"DD_SITE"`
	Title          string   `env:"TITLE,required"`
	Text           string   `env:"TEXT,required"`
	Tags           []string `env:"TAGS"`
	AggregationKey string   `env:"AGGREGATION_KEY"`
	RelatedEventId int64    `env:"RELATED_EVENT_ID"`
}

type output struct {
	ID int64 `json:"id"`
	step.Outputs
}

func (s *EventPost) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventPost) Run() (int, []byte, error) {
	//prepare request body
	body := datadog.EventCreateRequest{
		Title:          s.Title,
		Text:           s.Text,
		Tags:           &s.Tags,
		AggregationKey: &s.AggregationKey,
		RelatedEventId: &s.RelatedEventId,
		SourceTypeName: datadog.PtrString(stackpulseSource),
	}

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()

	// send request
	apiClient := datadog.NewAPIClient(configuration)
	resp, r, err := apiClient.EventsApi.CreateEvent(ctx).Body(body).Execute()
	if err != nil {
		apiError, ok := err.(datadog.GenericOpenAPIError)
		if ok && r != nil && r.StatusCode != http.StatusAccepted {
			return step.ExitCodeFailure, nil, fmt.Errorf("failed to create event. got response code: %d %s",
				r.StatusCode, apiError.Body())
		} else {
			return step.ExitCodeFailure, nil, err
		}
	}

	defer r.Body.Close()
	respBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare output
	var eventID int64
	if resp.Id != nil {
		eventID = *resp.Id
	}

	stepOutput := output{
		ID:      eventID,
		Outputs: step.Outputs{Object: string(respBody)},
	}

	jsonOutput, err := json.Marshal(stepOutput)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&EventPost{})
}
