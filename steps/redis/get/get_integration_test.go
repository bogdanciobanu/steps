// +build integration

package main

import (
	"testing"

	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps-sdk-go/testutil/container"
	"github.com/stackpulse/steps-sdk-go/testutil/container/matcher"
	"github.com/stackpulse/steps/redis/base/redistest"
)

func TestRedisGet_Run(t *testing.T) {
	testImage := container.CurrentStepImage(t)
	serviceUrls := redistest.SetupRedis(t)

	cases := []container.Test{
		{
			Name:         "no params",
			Image:        testImage,
			Envs:         map[string]string{},
			Args:         []string{},
			WantExitCode: step.ExitCodeFailure,
			WantError:    "failed init step arguments, env: required environment variable",
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
			Name:         "key not found",
			Image:        testImage,
			Envs:         map[string]string{"KEY": "not-exist", "REDIS_URL": serviceUrls.FullURL},
			Args:         []string{},
			WantExitCode: step.ExitCodeFailure,
			WantError:    "key not found",
			WantOutput:   nil,
		},
		{
			Name:         "numerical key - found",
			Image:        testImage,
			Envs:         map[string]string{"KEY": "number", "REDIS_URL": serviceUrls.FullURL},
			Args:         []string{},
			WantExitCode: step.ExitCodeOK,
			WantError:    "",
			WantOutput:   matcher.StepOutputTextEqual("42"),
		},
		{
			Name:         "numerical key - string",
			Image:        testImage,
			Envs:         map[string]string{"KEY": "string", "REDIS_URL": serviceUrls.FullURL},
			Args:         []string{},
			WantExitCode: step.ExitCodeOK,
			WantError:    "",
			WantOutput:   matcher.StepOutputTextEqual("foo"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, tc.Run)
	}
}
