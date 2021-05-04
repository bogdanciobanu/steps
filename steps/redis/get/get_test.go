package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/redis"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func ExecuteStep(envs map[string]string, stepImagePath string) (string, error) {
	req := testcontainers.ContainerRequest{
		Image:       stepImagePath,
		Env:         envs,
		NetworkMode: "host",
		SkipReaper:  true,
	}

	stepC, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return "", err
	}

	time.Sleep(1 * time.Second)

	r, err := stepC.Logs(context.Background())
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ParseStepOutput(stepOutput string) (map[string]string, error) {
	if !strings.Contains(stepOutput, env.EndMarkerDefault) {
		return nil, fmt.Errorf("failed to parse step output, no end marker found in output")
	}

	parsedOutput := strings.Split(stepOutput, env.EndMarkerDefault)

	jsonOutput := parsedOutput[1]

	var unmarshalJsonOutput map[string]string

	err := json.Unmarshal([]byte(jsonOutput), &unmarshalJsonOutput)

	return unmarshalJsonOutput, err
}

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
	container, err := gnomock.Start(p)
	if err != nil {
		assert.Fail(t, "failed to create redis container: %w", err)
	}

	return NewServiceUrls("redis://", container.Host, strconv.Itoa(container.DefaultPort()))
}

func TestStep(t *testing.T) {
	serviceUrls := SetupRedis(t)

	tests := []struct {
		name           string
		envs           map[string]string
		shouldError    bool
		errorContains  string
		expectedOutput string
	}{
		{
			"no params",
			map[string]string{},
			true,
			"failed init step arguments, env: required environment variable",
			"",
		},
		{
			"failed to connect redis - invalid dns",
			map[string]string{"KEY": "mykey", "REDIS_URL": "redis://invalid-hostname"},
			true,
			"no such host",
			"",
		},
		{
			"failed to connect redis - invalid ip",
			map[string]string{"KEY": "mykey", "REDIS_URL": "redis://127.0.0.2"},
			true,
			"connection refused",
			"",
		},
		{
			"key not found",
			map[string]string{"KEY": "not-exist", "REDIS_URL": serviceUrls.FullUrl},
			true,
			"key not found",
			"",
		},
		{
			"numerical key - found",
			map[string]string{"KEY": "number", "REDIS_URL": serviceUrls.FullUrl},
			false,
			"",
			"42",
		},
		{
			"numerical key - string",
			map[string]string{"KEY": "string", "REDIS_URL": serviceUrls.FullUrl},
			false,
			"",
			"foo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stepOutput, err := ExecuteStep(test.envs, "us-docker.pkg.dev/stackpulse/public/steps/redis/get")
			if err != nil {
				assert.Fail(t, "failed to execute the step: %w", err)
			}

			if test.shouldError {
				assert.Contains(t, stepOutput, test.errorContains)

			} else {
				parsedStepOutput, err := ParseStepOutput(stepOutput)
				if err != nil {
					assert.Fail(t, "failed to parse step output: %w", err)
				}

				assert.Equal(t, test.expectedOutput, parsedStepOutput["output"])
			}
		})
	}
}
