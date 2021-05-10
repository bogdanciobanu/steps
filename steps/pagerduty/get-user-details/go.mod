module github.com/stackpulse/steps/steps/pagerduty/get-user-details

go 1.16

replace github.com/stackpulse/steps/pagerduty/base v0.0.0 => ../base

require (
	github.com/PagerDuty/go-pagerduty v1.3.0 // indirect
	github.com/stackpulse/steps-sdk-go v0.0.0-20210329111118-cc87e9772586 // indirect
	github.com/stackpulse/steps/pagerduty/base v0.0.0
)
