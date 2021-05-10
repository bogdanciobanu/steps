package main

import (
	"encoding/json"

	pd "github.com/PagerDuty/go-pagerduty"
	envconf "github.com/caarlos0/env/v6"

	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/pagerduty/base"
)

type Args struct {
	base.Args
	IsEarliest            bool     `env:"EARLIEST" envDefault:"false"`
	EscalationPolicyIDs   []string `env:"ESCALATION_POLICIES_IDS" envSeparator:","`
	EscalationPolicyNames []string `env:"ESCALATION_POLICIES_NAMES" envSeparator:","`
	TimeZone              string   `env:"TIME_ZONE"`
	Since                 string   `env:"SINCE"`
	Until                 string   `env:"UNTIL"`
}

type PagerDutyOnCallsList struct {
	args   Args
	client *pd.Client
}

func (p *PagerDutyOnCallsList) Init() error {
	err := envconf.Parse(&p.args)
	if err != nil {
		return err
	}

	p.client = pd.NewClient(p.args.PdToken)
	return nil
}

func (p *PagerDutyOnCallsList) buildListOpts() (pd.ListOnCallOptions, error) {
	onCallOpts := pd.ListOnCallOptions{}

	if p.args.IsEarliest {
		onCallOpts.Earliest = p.args.IsEarliest
		log.Debug("Listing only earliest on calls")
	}

	if len(p.args.EscalationPolicyNames) != 0 {
		idsByNames, err := p.getEscalationPolicyIdsByNames(p.args.EscalationPolicyNames)
		if err != nil {
			return onCallOpts, err
		}
		p.args.EscalationPolicyIDs = append(p.args.EscalationPolicyIDs, idsByNames...)
		log.Debug("Listing on calls for escaltion policies %v", p.args.EscalationPolicyNames)
	}

	if len(p.args.EscalationPolicyIDs) != 0 {
		onCallOpts.EscalationPolicyIDs = p.args.EscalationPolicyIDs
		log.Debug("Listing on calls for escaltion policies %v", p.args.EscalationPolicyIDs)
	}

	if p.args.Since != "" {
		onCallOpts.Since = p.args.Since
		log.Debug("Listing on calls since %v", p.args.Since)
	}

	if p.args.Until != "" {
		onCallOpts.Until = p.args.Until
		log.Debug("Listing on calls until %v", p.args.Until)
	}

	if p.args.TimeZone != "" {
		onCallOpts.TimeZone = p.args.TimeZone
		log.Debug("Listing on calls with time zone %v", p.args.TimeZone)
	}

	return onCallOpts, nil
}

func (p *PagerDutyOnCallsList) getEscalationPolicyIdsByNames(epNames []string) (epIDs []string, err error) {
	epResp, err := p.client.ListEscalationPolicies(pd.ListEscalationPoliciesOptions{})
	if err != nil {
		return nil, err
	}

	for _, ep := range epResp.EscalationPolicies {
		if base.Contains(epNames, ep.Name) {
			epIDs = append(epIDs, ep.ID)
		}
		if len(epIDs) >= len(epNames) {
			break
		}
	}
	return epIDs, nil
}

func (p *PagerDutyOnCallsList) Run() (exitCode int, output []byte, err error) {
	onCallOpts, err := p.buildListOpts()
	getUserOpts := pd.GetUserOptions{}
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	res, err := p.client.ListOnCalls(onCallOpts)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	for key, value := range res.OnCalls {
		userId := value.User.ID
		user, err := p.client.GetUser(userId, getUserOpts)
		if err != nil {
			return step.ExitCodeFailure, nil, err
		}
		res.OnCalls[key].User.Email = user.Email
	}
	marshaledOnCall, err := json.Marshal(res.OnCalls)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	return step.ExitCodeOK, marshaledOnCall, nil
}

func main() {
	step.Run(&PagerDutyOnCallsList{})
}
