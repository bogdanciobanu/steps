package main

import (
	"encoding/json"
	"fmt"

	pd "github.com/PagerDuty/go-pagerduty"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/pagerduty/base"
)

type Args struct {
	base.Args
	UserId string `env:"USER_ID,required"`
}

type PagerDutyGetUser struct {
	args   Args
	client *pd.Client
}

func (p *PagerDutyGetUser) Init() error {
	if err := env.Parse(&p.args); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}
	if len(p.args.UserId) == 0 {
		return fmt.Errorf("no user id provided")
	}
	p.client = pd.NewClient(p.args.PdToken)
	return nil
}

func (p *PagerDutyGetUser) Run() (int, []byte, error) {
	var ret []byte
	getUserOpts := pd.GetUserOptions{}
	res, err := p.client.GetUser(p.args.UserId, getUserOpts)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	if ret, err = json.Marshal(res); err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json marshal: %w", err)
	}
	return step.ExitCodeOK, ret, nil
}

func main() {
	step.Run(&PagerDutyGetUser{})
}
