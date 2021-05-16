module github.com/stackpulse/steps/redis/client-list

go 1.14

require (
	github.com/caarlos0/env/v6 v6.5.0
	github.com/go-redis/redis/v8 v8.4.4
	github.com/stackpulse/steps-sdk-go v0.0.0-20210512171627-c9c9d34819ce
	github.com/stackpulse/steps/redis/base v0.0.0
)

replace github.com/stackpulse/steps/redis/base v0.0.0 => ../base
