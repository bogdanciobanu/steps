package main

import (
	"context"
	"fmt"

	"github.com/Jeffail/gabs/v2"
	envconf "github.com/caarlos0/env/v6"
	"github.com/go-redis/redis/v8"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/redis/base"
)

type Args struct {
	base.Args
	Key string `env:"KEY,required"`
}

type RedisGet struct {
	args        Args
	redisClient *redis.Client
}

func (l *RedisGet) Init() error {
	err := envconf.Parse(&l.args)
	if err != nil {
		return err
	}

	redisClient, err := base.InitRedisClient(l.args.Args)
	if err != nil {
		return err
	}

	l.redisClient = redisClient

	return nil
}

func (l *RedisGet) Run() (exitCode int, output []byte, err error) {
	gc := gabs.New()

	cmd := l.redisClient.Get(context.Background(), l.args.Key)
	if cmd == nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("invalid result returned from get operation")
	}

	val, err := cmd.Result()
	if err != nil {
		if err == redis.Nil {
			return step.ExitCodeFailure, nil, fmt.Errorf("key not found")
		}
		return step.ExitCodeFailure, nil, fmt.Errorf("failed getting key value: validate key is valid: %w", err)
	}

	if _, err = gc.Set(val, "output"); err != nil {
		return step.ExitCodeFailure, []byte(val), err
	}

	return step.ExitCodeOK, gc.Bytes(), nil
}

func main() {
	step.Run(&RedisGet{})
}
