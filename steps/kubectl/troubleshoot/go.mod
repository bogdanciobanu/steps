module github.com/stackpulse/steps/steps/kubectl/troubleshoot

go 1.14

require (
	github.com/olekukonko/tablewriter v0.0.5
	github.com/stackpulse/steps-sdk-go v0.0.0-20210418111602-8f22092d9c5a
	github.com/stackpulse/steps/kubectl/base v0.0.0
)

replace github.com/stackpulse/steps/kubectl/base v0.0.0 => ../base
