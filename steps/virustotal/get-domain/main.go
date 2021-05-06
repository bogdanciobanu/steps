package main

import (
	"encoding/json"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"

	vt "github.com/VirusTotal/vt-go"
)

type VirustotalGetUrl struct {
	VtApiKey string `env:"VT_API_KEY,required"`
	Domain   string `env:"DOMAIN,required"`
}

type stepOutput struct {
	Reputation int         `json:"reputation"`
	Stats      statsOutput `json:"stats"`
	step.Outputs
}

type statsOutput struct {
	Harmless   int `json:"harmless"`
	Malicious  int `json:"malicious"`
	Suspicious int `json:"suspicious"`
	Timeout    int `json:"timeout"`
	Undetected int `json:"undetected"`
}

type apiResponse struct {
	Attributes struct {
		LastAnalysisStats statsOutput `json:"last_analysis_stats"`
		AsOwner           string      `json:"as_owner"`
		Reputation        int         `json:"reputation"`
	} `json:"attributes"`
}

func (s *VirustotalGetUrl) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *VirustotalGetUrl) Run() (exitCode int, output []byte, err error) {
	//prepare client
	client := vt.NewClient(s.VtApiKey)

	//send request
	rawResp, err := client.Get(vt.URL("domains/%s", s.Domain))
	if err != nil {
		//err is present if invalid JSON or non 2xx status code
		return step.ExitCodeFailure, nil, err
	}

	//parse response
	rawRes, err := rawResp.Data.MarshalJSON()
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	var resp apiResponse
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//unmarshal full response
	var fullResp map[string]interface{}
	err = json.Unmarshal(rawRes, &fullResp)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare output
	out := stepOutput{
		Reputation: resp.Attributes.Reputation,
		Stats:      resp.Attributes.LastAnalysisStats,
		Outputs:    step.Outputs{Object: fullResp},
	}
	outputBytes, err := json.Marshal(out)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//success
	return step.ExitCodeOK, outputBytes, nil
}

func main() {
	step.Run(&VirustotalGetUrl{})
}
