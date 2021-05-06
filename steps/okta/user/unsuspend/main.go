package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"

	"github.com/okta/okta-sdk-golang/okta"
)

type UserUnsuspend struct {
	OktaApiToken string `env:"OKTA_API_TOKEN,required"`
	OktaDomain   string `env:"OKTA_DOMAIN,required"`
	UserId       string `env:"USER_ID,required"`
	ctx          context.Context
}

type stepResult string

const (
	stepResultSuccess      stepResult = "success"
	stepResultInvalidID    stepResult = "invalid user id"
	stepResultInvalidState stepResult = "invalid user state"
)

func (s *UserUnsuspend) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	//default context
	s.ctx = context.Background()

	return nil
}

func (s *UserUnsuspend) Run() (int, []byte, error) {
	//create client
	oktaClient, err := okta.NewClient(s.ctx, okta.WithOrgUrl(fmt.Sprintf("https://%s", s.OktaDomain)), okta.WithToken(s.OktaApiToken))
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//send request
	_, resp, err := oktaClient.User.ActivateUser(s.UserId, nil)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return step.ExitCodeOK, []byte(stepResultInvalidID), nil
		} else if resp.StatusCode == http.StatusForbidden {
			return step.ExitCodeOK, []byte(stepResultInvalidState), nil
		}
		return step.ExitCodeFailure, nil, err
	}
	return step.ExitCodeOK, []byte(stepResultSuccess), nil
}

func main() {
	step.Run(&UserUnsuspend{})
}
