module github.com/stackpulse/steps/redis/get

go 1.14

require (
	github.com/Jeffail/gabs/v2 v2.6.0
	github.com/caarlos0/env/v6 v6.5.0
	github.com/go-redis/redis/v8 v8.4.4
	github.com/orlangure/gnomock v0.14.0
	github.com/stackpulse/steps-sdk-go v0.0.0-20210511115437-47cffa89c1f2
	github.com/stackpulse/steps/redis/base v0.0.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/stackpulse/steps/redis/base v0.0.0 => ../base
