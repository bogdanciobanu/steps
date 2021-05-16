package redistest

import (
	"strconv"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/redis"
	"github.com/stackpulse/steps-sdk-go/testutil/container"
)

func SetupRedis(t *testing.T) container.ServiceURL {
	t.Helper()
	vs := make(map[string]interface{})

	vs["string"] = "foo"
	vs["number"] = 42
	vs["boolean"] = true

	p := redis.Preset(redis.WithValues(vs))
	redisContainer, err := gnomock.Start(p)
	if err != nil {
		t.Fatalf("failed to create redis redisContainer: %w", err)
	}

	return container.NewServiceURL("redis://", redisContainer.Host, strconv.Itoa(redisContainer.DefaultPort()))
}
