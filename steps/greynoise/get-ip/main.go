package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
)

type GreynoiseGetIp struct {
	Ip       string `env:"IP,required"`
	GnApiKey string `env:"GN_API_KEY"`
}

type stepOutput struct {
	Name           string `json:"name"`
	Noise          bool   `json:"noise"`
	Riot           bool   `json:"riot"`
	Classification string `json:"classification"`
	step.Outputs
}

type apiResponse struct {
	Name           string `json:"name"`
	Noise          bool   `json:"noise"`
	Riot           bool   `json:"riot"`
	Classification string `json:"classification"`
}

func (s *GreynoiseGetIp) Init() error {
	err := env.Parse(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *GreynoiseGetIp) Run() (int, []byte, error) {
	//prepare request
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.greynoise.io/v3/community/%s", s.Ip), nil)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}
	req.Header.Set("key", s.GnApiKey)

	//send request
	rawResp, err := http.DefaultClient.Do(req)
	if err != nil || rawResp.StatusCode/100 != 2 {
		//parse response
		defer rawResp.Body.Close()
		respBytes, err := ioutil.ReadAll(rawResp.Body)
		if err != nil {
			return step.ExitCodeFailure, nil, err
		}
		return step.ExitCodeFailure, respBytes, err
	}

	//parse response
	defer rawResp.Body.Close()
	respBytes, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	var resp apiResponse
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		internalErr := fmt.Errorf("could not unmarshal response to JSON: %w; res: %s", err, string(respBytes))
		return step.ExitCodeFailure, nil, internalErr
	}

	var fullResp map[string]interface{}
	err = json.Unmarshal(respBytes, &fullResp)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//prepare step output
	output := stepOutput{
		Name:           resp.Name,
		Noise:          resp.Noise,
		Riot:           resp.Riot,
		Classification: resp.Classification,
		Outputs:        step.Outputs{Object: fullResp},
	}
	outputBytes, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	//success
	return step.ExitCodeOK, outputBytes, nil
}

func main() {
	step.Run(&GreynoiseGetIp{})
}
