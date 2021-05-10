package main

import (
	"context"
	"encoding/json"
	"fmt"

	github "github.com/shurcooL/githubv4"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/github/base"
)

type Args struct {
	base.Args
	Repository       string `env:"REPOSITORY,required"`
	RepositoryNodeID string
	BaseRefName      string `env:"BASE_BRANCH" envDefault:"main"`
	HeadRefName      string `env:"HEAD_BRANCH,required"`
	Title            string `env:"TITLE,required"`
}

type GHCreatePullRequest struct {
	*base.GithubClient
	args Args
}

type Outputs struct {
	step.Outputs
	Branch string `json:"branch"`
	Link   string `json:"link"`
	Author string `json:"author"`
	State  string `json:"state"`
}

func (gh *GHCreatePullRequest) Init() error {
	var args Args
	baseClient, err := base.NewGithubClient(&args)
	if err != nil {
		return err
	}

	gh.GithubClient = baseClient
	gh.args = args

	return nil
}

func (gh *GHCreatePullRequest) createMutationInput() github.CreatePullRequestInput {
	return github.CreatePullRequestInput{
		RepositoryID: github.String(gh.args.RepositoryNodeID),
		BaseRefName:  github.String(gh.args.BaseRefName),
		HeadRefName:  github.String(gh.args.HeadRefName),
		Title:        github.String(gh.args.Title),
	}
}

func (gh *GHCreatePullRequest) getRepositoryNodeID() (string, error) {
	query := &repositoryIDQuery
	gh.QueryVariables["repository"] = github.String(gh.args.Repository)

	err := gh.Client.Query(context.Background(), &query, gh.QueryVariables)
	if err != nil {
		return "", err
	}

	return repositoryIDQuery.Repository.ID, nil
}

func (gh *GHCreatePullRequest) Run() (int, []byte, error) {
	var err error
	output := Outputs{}
	mutation := &createPullRequestMutation

	gh.args.RepositoryNodeID, err = gh.getRepositoryNodeID()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create-pull-request failed with error: error: %w", err)
	}
	input := gh.createMutationInput()

	err = gh.Client.Mutate(context.Background(), mutation, input, nil)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create-pull-request failed with error: error: %w", err)
	}

	pathToPayload := mutation.CreatePullRequest.PullRequest
	output.Link = pathToPayload.Permalink
	output.Branch = pathToPayload.HeadRefName
	output.State = pathToPayload.State
	output.Author = pathToPayload.Author.Login

	marshalled, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create-pull-request failed with error: error: %w", err)
	}
	return step.ExitCodeOK, marshalled, err
}

func main() {
	step.Run(&GHCreatePullRequest{})
}

var createPullRequestMutation struct {
	CreatePullRequest struct {
		PullRequest struct {
			HeadRefName string
			Permalink   string
			Author      struct {
				Login string
			}
			State string
		}
	} `graphql:"createPullRequest(input: $input)"`
}

var repositoryIDQuery struct {
	Repository struct {
		ID string
	} `graphql:"repository(name: $repository, owner: $owner)"`
}
