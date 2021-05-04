module github.com/stackpulse/steps/redis/get

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/go-redis/redis/v8 v8.4.4
	github.com/google/uuid v1.2.0
	github.com/orlangure/gnomock v0.14.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210314133745-61086c27983f
	github.com/stackpulse/steps/redis/base v0.0.0
	github.com/stretchr/testify v1.7.0
	github.com/testcontainers/testcontainers-go v0.10.0
)

replace github.com/stackpulse/steps/redis/base v0.0.0 => ../base
