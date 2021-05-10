package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/redis"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps-sdk-go/stesting/container"
	"github.com/stretchr/testify/assert"
)

type ServiceUrls struct {
	Schema      string
	Host        string
	Port        string
	HostAndPort string
	FullUrl     string
}

func NewServiceUrls(schema, host, port string) ServiceUrls {
	if runtime.GOOS == "darwin" {
		host = strings.Replace(host, "127.0.0.1", "docker.for.mac.localhost", 1)
	}

	return ServiceUrls{
		Schema:      schema,
		Host:        host,
		Port:        port,
		HostAndPort: fmt.Sprintf("%s:%s", host, port),
		FullUrl:     fmt.Sprintf("%s%s:%s", schema, host, port),
	}
}

func SetupRedis(t *testing.T) ServiceUrls {
	vs := make(map[string]interface{})

	vs["string"] = "foo"
	vs["number"] = 42
	vs["boolean"] = true

	p := redis.Preset(redis.WithValues(vs))
	redisContainer, err := gnomock.Start(p)
	if err != nil {
		assert.Fail(t, "failed to create redis redisContainer: %w", err)
	}

	return NewServiceUrls("redis://", redisContainer.Host, strconv.Itoa(redisContainer.DefaultPort()))
}

func TestRedisListClientList_Run(t *testing.T) {
	serviceUrls := SetupRedis(t)

	cases := []container.TestCase{
		{
			Name:			"no params",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/client-list",
			Envs:			map[string]string{},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError: 		"failed init step arguments, env: required environment variable",
			WantOutput: 	"",
		},
		{
			Name:			"failed to connect redis - invalid dns",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/client-list",
			Envs:			map[string]string{"KEY": "mykey", "REDIS_URL": "redis://invalid-hostname"},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError:		"no such host",
			WantOutput:		"",
		},
		{
			Name:			"failed to connect redis - invalid ip",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/client-list",
			Envs:			map[string]string{"KEY": "mykey", "REDIS_URL": "redis://127.0.0.2"},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError:		"connection refused",
			WantOutput:		"",
		},
		{
			Name:			"negative limit",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/client-list",
			Envs:			map[string]string{"LIMIT": "-1", "REDIS_URL": serviceUrls.FullUrl},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeFailure,
			WantError:		"LIMIT cannot be a negative number",
			WantOutput:		"",
		},
		{
			Name:			"list one client",
			Image:			"us-docker.pkg.dev/stackpulse/public/redis/client-list",
			Envs:			map[string]string{"REDIS_URL": serviceUrls.FullUrl},
			Args: 			[]string{},
			WantExitCode: 	step.ExitCodeOK,
			WantError:		"",
			WantOutput:		"",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.Run(t)
		})
	}
}
