package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
)

type MessageSend struct {
	Webhook   string `env:"WEBHOOK,required"`
	Message   string `env:"MESSAGE,required"`
	ThreadKey string `env:"THREAD_KEY"`

	parsedWebhook *url.URL
}

type output struct {
	ThreadName string `json:"thread_name"`
	SpaceName  string `json:"space_name"`
	step.Outputs
}

type ChatJson struct {
	Name         string        `json:"name"`
	Sender       Sender        `json:"sender"`
	Text         string        `json:"text"`
	Cards        []interface{} `json:"cards"`
	PreviewText  string        `json:"previewText"`
	Annotations  []interface{} `json:"annotations"`
	Thread       Thread        `json:"thread"`
	Space        Space         `json:"space"`
	FallbackText string        `json:"fallbackText"`
	ArgumentText string        `json:"argumentText"`
	Attachment   []interface{} `json:"attachment"`
	CreateTime   time.Time     `json:"createTime"`
}
type Sender struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Email       string `json:"email"`
	DomainID    string `json:"domainId"`
	Type        string `json:"type"`
	IsAnonymous bool   `json:"isAnonymous"`
}
type Thread struct {
	Name string `json:"name"`
}
type Space struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	SingleUserBotDm bool   `json:"singleUserBotDm"`
	Threaded        bool   `json:"threaded"`
	DisplayName     string `json:"displayName"`
}

func (s *MessageSend) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	parsedWebhook, err := url.Parse(s.Webhook)
	if err != nil {
		return fmt.Errorf("parse webhook URL: %w", err)
	}
	s.parsedWebhook = parsedWebhook
	return nil
}

func (s *MessageSend) Run() (int, []byte, error) {
	if s.ThreadKey != "" {
		q := s.parsedWebhook.Query()
		q.Set("threadKey", s.ThreadKey)
		s.parsedWebhook.RawQuery = q.Encode()
	}
	// Generated post body
	postBody, _ := json.Marshal(map[string]string{
		"text": s.Message,
	})

	requestBody := bytes.NewBuffer(postBody)

	webhookURL := s.parsedWebhook.String()
	log.Debugln("Sending post request to: %s. Body: %s", webhookURL, string(postBody))
	// Send post to webhook
	resp, err := http.Post(webhookURL, "application/json", requestBody)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get webhook : %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("read body : %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return step.ExitCodeFailure, body, fmt.Errorf("got non 20X HTTP response code: %d(%s)", resp.StatusCode, resp.Status)
	}

	// Unmarshal response body
	var jsonResp ChatJson
	err = json.Unmarshal(body, &jsonResp)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json unmarshal (%s) : %w", body, err)
	}

	// Generate step output
	out := output{
		ThreadName: jsonResp.Thread.Name,
		SpaceName:  jsonResp.Space.Name,
		Outputs:    step.Outputs{Object: jsonResp},
	}

	jsonOutput, err := json.Marshal(&out)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("marshal output: %w", err)
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&MessageSend{})
}
