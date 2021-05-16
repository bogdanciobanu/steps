// +build integration

package main

import (
	"testing"

	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps-sdk-go/testutil/container"
	"github.com/stackpulse/steps/redis/base/redistest"
)

func TestRedisListClientList_Run(t *testing.T) {
	testImage := container.CurrentStepImage(t)
	serviceURL := redistest.SetupRedis(t)

	cases := []container.Test{
		{
			Name:         "no params",
			Image:        testImage,
			Envs:         map[string]string{},
			Args:         []string{},
			WantExitCode: step.ExitCodeFailure,
			WantError:    "invalid arguments",
			WantOutput:   nil,
		},
		{
			Name:         "failed to connect redis - invalid dns",
			Image:        testImage,
			Envs:         map[string]string{"KEY": "mykey", "REDIS_URL": "redis://invalid-hostname.internal"},
			Args:         []string{},
			WantExitCode: step.ExitCodeFailure,
			WantError:    "no such host",
			WantOutput:   nil,
		},
		{
			Name:         "failed to connect redis - invalid ip",
			Image:        testImage,
			Envs:         map[string]string{"KEY": "mykey", "REDIS_URL": "redis://127.0.0.2"},
			Args:         []string{},
			WantExitCode: step.ExitCodeFailure,
			WantError:    "connection refused",
			WantOutput:   nil,
		},
		{
			Name:         "negative limit",
			Image:        testImage,
			Envs:         map[string]string{"LIMIT": "-2", "REDIS_URL": serviceURL.FullURL},
			Args:         []string{},
			WantExitCode: step.ExitCodeFailure,
			WantError:    "LIMIT cannot be a negative number",
			WantOutput:   nil,
		},
		{
			Name:         "list one client",
			Image:        testImage,
			Envs:         map[string]string{"REDIS_URL": serviceURL.FullURL},
			Args:         []string{},
			WantExitCode: step.ExitCodeOK,
			WantError:    "",
			WantOutput:   nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, tc.Run)
	}
}
