module github.com/stackpulse/steps/steps/github/get-file

go 1.15

require (
	github.com/shurcooL/githubv4 v0.0.0-20201206200315-234843c633fa
	github.com/shurcooL/graphql v0.0.0-20200928012149-18c5c3165e3a // indirect
	github.com/stackpulse/steps-sdk-go v0.0.0-20210329111118-cc87e9772586
	github.com/stackpulse/steps/github/base v0.0.0
)

replace github.com/stackpulse/steps/github/base v0.0.0 => ../base
