package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps-sdk-go/upload"
)

const (
	NetworkTimeout = 15 * time.Second
)

type Args struct {
	URL     string `env:"URL,required"`
	API_KEY string `env:"API_KEY", envDefault:""`
}

type GrafanaGetPanel struct {
	args Args
}

func (l *GrafanaGetPanel) Init() error {
	err := env.Parse(&l.args)
	if err != nil {
		return err
	}

	return nil
}

func (l *GrafanaGetPanel) Run() (int, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), NetworkTimeout)
	defer cancel()
	client := &http.Client{}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, l.args.URL, nil)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("newRequest: %w", err)
	}

	if l.args.API_KEY != "" {
		// add authorization header with api key to the req
		var bearer = "Bearer " + l.args.API_KEY
		request.Header.Add("Authorization", bearer)
	}

	resp, err := client.Do(request)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return step.ExitCodeFailure, nil, fmt.Errorf("got non 20X HTTP response code: %d(%s)", resp.StatusCode, resp.Status)
	}

	output, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("readPanel: %w", err)
	}

	if err = upload.RichOutput(ctx, ioutil.NopCloser(bytes.NewBuffer(output)), upload.ContentTypePNG); err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("uploadPanel: %w", err)
	}

	result := []byte(base64.StdEncoding.EncodeToString(output))
	return step.ExitCodeOK, result, nil
}

func main() {
	step.Run(&GrafanaGetPanel{})
}
