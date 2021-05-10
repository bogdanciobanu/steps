module github.com/stackpulse/steps/steps/kubectl/heapdump

go 1.16

require (
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/stackpulse/steps-sdk-go v0.0.0-20210413164812-58fc219b8518
	github.com/stackpulse/steps/kubectl/base v0.0.0
)

replace github.com/stackpulse/steps/kubectl/base v0.0.0 => ../base
