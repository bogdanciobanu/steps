package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

type UserSuspend struct {
	OktaApiToken string `env:"OKTA_API_TOKEN,required"`
	OktaDomain   string `env:"OKTA_DOMAIN,required"`
	UserId       string `env:"USER_ID,required"`
}

func (s *UserSuspend) Init() *step.Error {
	err := env.Parse(s)
	if err != nil {
		return &step.Error{
			Code:    step.InvalidArgument,
			Message: err.Error(),
		}
	}

	return nil
}

func (s *UserSuspend) Run() (step.Output, *step.Error) {
	// create client
	ctx, oktaClient, err := okta.NewClient(context.Background(), okta.WithOrgUrl(fmt.Sprintf("https://%s", s.OktaDomain)), okta.WithToken(s.OktaApiToken))
	if err != nil {
		return nil, &step.Error{
			Code:    step.Internal,
			Message: "failed creating the okta sdk client",
			Verbose: err.Error(),
		}
	}

	// send request
	resp, err := oktaClient.User.SuspendUser(ctx, s.UserId)
	stepError := step.Error{}
	if err != nil {
		switch {
		case resp.StatusCode == http.StatusNotFound:
			stepError.Code = step.NotFound
			stepError.Message = "invalid user id: user id doesn't exist"
		case resp.StatusCode == http.StatusBadRequest:
			stepError.Code = step.InvalidArgument
			stepError.Message = "invalid user state: user must be in the ACTIVE state"
		case resp.StatusCode/100 != 2:
			stepError.Code = step.Unavailable
			stepError.Message = fmt.Sprintf("request returned non success status: %d", resp.StatusCode)
		}
	}

	return nil, &stepError
}

func main() {
	step.Run(&UserSuspend{})
}
