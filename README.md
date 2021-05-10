# ![stackpulse-logo](vendors/stackpulse.svg) StackPulse Steps

[![Build][badge_ci]][link_circle]
[![Contributors][contributors-shield]][contributors-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][badge_issues]][issues-url]
[![License][badge_license]][link_license]

This repository is the official library of StackPulse steps.

StackPulse steps are containerized applications that can be composed together to form a playbook. Steps can perform any arbitrary computational task but are most often calls to different services and APIs, allowing the playbook composer to interact with their stack in an automated way. This enables faster and more streamlined handling of alerts and remediation of incidents.

Just like any other container, steps can receive input in the form of run arguments or environment variables and have output directed to stdout and stderr. The standard form for a step to issue its output is as a JSON object as defined in the [StackPulse steps Golang SDK](https://github.com/stackpulse/steps-sdk-go).

## Development

### Build tool

The build tool used by the steps repository is called baur. You must install it on your machine from [here](https://github.com/simplesurance/baur/releases).  The currently supported version we use is *1.0-rc2*.

Baur requires PostgreSQL to operate. After building a step once, the step's file hashes will be stored into postgres. Subsequent attempts to build the step will not do anything unless one of the files the step depends on are changed. 

### Building all steps locally

In order to start postgres docker container on you must run `make pg`. This will start a docker container on your machine listening on port `15432`. Baur is configured to use postgres running on that port to store its data.

You can build all steps by executing `make local` command. This process takes a while. Once all steps are built you will enjoy faster subsequent builds.

### Building a specific step locally

If you want to build a single step you can run `baur run step/name`, for example `baur run redis/get`. If the step or its dependencies were changed the step will be rebuilt.
The result of a successful build is a docker image with a tag named of the current branch (i.e: `us-docker.pkg.dev/stackpulse/public/redis/keys:MyNewStep`)


### Step directory structure

Step normally belongs to a family, a family (i.e `redis`) will contain several steps that are related to it. It will normally be grouped by some common service all steps operate on.

Often times, there will be shared code for all the steps. It can either be a base Docker image that all steps inherit from, some common Go code, etc. 
```text
./steps/redis
|-- base
|   |-- base.go
|   |-- go.mod
|   `-- go.sum
|-- bigkeys
|   |-- Dockerfile
|   |-- big-keys.py
|   `-- manifest.yaml
|-- client-list
|   |-- Dockerfile
|   |-- go.mod
|   |-- go.sum
|   |-- main.go
|   `-- manifest.yaml
```
Any change to `base` will result all sibling steps to be rebuilt. A change to a specific step (`get`) will rebuild only that particular step.

### Creating a new step
1. Decide to which family the step belongs, does such family already exist? Check the `./steps/` directory
2. If such directory already exists, review steps currently in that family. Aknowledge the `base` directory and see what is already implemented there.
3. Create a your step under `./steps/family/`, it should contain a `Dockerfile`.
4. Once your step contains `Dockerfile` run `make apps` to generate baur application file for the newly created steps, the file created will be called `.app.toml`
5. After running the command you should be able to run `baur run family/my-new-step` to build it.

## Contributing

Pull requests adding new steps and submitting bug fixes are welcome.

Please note that all new steps must use the [StackPulse steps Golang SDK](https://github.com/stackpulse/steps-sdk-go) and be submitted with a step manifest containing input and output examples.

In the case of bug fixes, where appropriate, please provide an accompanying unit test that validates and verifies the proposed fix.

## Troubleshooting

Please feel free to open an issue for any problem or you may encounter while using one of the steps in the repo. We welcome requests and suggestions for new steps as well.

## License

Distributed under the BSD-3-Clause License. See `LICENSE` for more information.

[badge_ci]:https://circleci.com/gh/stackpulse/steps.svg?style=shield
[contributors-shield]: https://img.shields.io/github/contributors/stackpulse/steps.svg?style=flat-square&maxAge=30
[contributors-url]: https://github.com/stackpulse/steps/graphs/contributors
[badge_issues]:https://img.shields.io/github/issues/stackpulse/steps.svg?style=flat-square&maxAge=30
[issues-url]: https://github.com/stackpulse/steps/issues
[stars-shield]: https://img.shields.io/github/stars/stackpulse/steps.svg?style=flat-square&maxAge=30
[stars-url]: https://github.com/stackpulse/steps/stargazers
[badge_license]:https://img.shields.io/github/license/stackpulse/steps.svg?style=flat-square&maxAge=30
[link_license]:https://github.com/stackpulse/steps/blob/master/LICENSE
[link_circle]:https://circleci.com/gh/stackpulse/steps
