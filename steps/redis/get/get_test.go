package main

import (
	"strconv"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/redis"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps-sdk-go/testutil/container"
	"github.com/stackpulse/steps-sdk-go/testutil/container/matcher"
	"github.com/stretchr/testify/assert"
)


func SetupRedis(t *testing.T) container.ServiceUrls {
	vs := make(map[string]interface{})

	vs["string"] = "foo"
	vs["number"] = 42
	vs["boolean"] = true

	p := redis.Preset(redis.WithValues(vs))
	redisContainer, err := gnomock.Start(p)
	if err != nil {
		assert.Fail(t, "failed to create redis redisContainer: %w", err)
	}

	return container.NewServiceUrls("redis://", redisContainer.Host, strconv.Itoa(redisContainer.DefaultPort()))
}

func TestRedisGet_Run(t *testing.T) {
	serviceUrls := SetupRedis(t)

	cases := []container.Test{
		{
			Name:			"no params",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/get",
			Envs:			map[string]string{},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError: 		"failed init step arguments, env: required environment variable",
			WantOutput: 	nil,
		},
		{
			Name:			"failed to connect redis - invalid dns",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/get",
			Envs:			map[string]string{"KEY": "mykey", "REDIS_URL": "redis://invalid-hostname"},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError:		"no such host",
			WantOutput:		nil,
		},
		{
			Name:			"failed to connect redis - invalid ip",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/get",
			Envs:			map[string]string{"KEY": "mykey", "REDIS_URL": "redis://127.0.0.2"},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError:		"connection refused",
			WantOutput:		nil,
		},
		{
			Name:			"key not found",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/get",
			Envs:			map[string]string{"KEY": "not-exist", "REDIS_URL": serviceUrls.FullURL},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError:		"key not found",
			WantOutput:		nil,
		},
		{
			Name:			"numerical key - found",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/get",
			Envs:			map[string]string{"KEY": "number", "REDIS_URL": serviceUrls.FullURL},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeOK,
			WantError:		"",
			WantOutput:		matcher.Text("42"),
		},
		{
			Name:			"numerical key - string",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/get",
			Envs:			map[string]string{"KEY": "string", "REDIS_URL": serviceUrls.FullURL},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeOK,
			WantError:		"",
			WantOutput:		matcher.Text("foo"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, tc.Run)
	}
}
