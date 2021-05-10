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
	Path       string `env:"FILE_PATH,required"`
	Repository string `env:"REPOSITORY,required"`
	Branch     string `env:"BRANCH,required" envDefault:"main"`
}

type GHGetFile struct {
	*base.GithubClient
	args Args
}

type Outputs struct {
	step.Outputs
	FileContent string `json:"file_content"`
}

func (gh *GHGetFile) Init() error {
	var args Args
	baseClient, err := base.NewGithubClient(&args)
	if err != nil {
		return err
	}

	gh.GithubClient = baseClient
	gh.args = args
	gh.addQueryVariables()

	return nil
}

func (gh *GHGetFile) addQueryVariables() {
	gh.QueryVariables["filePath"] = github.String(gh.args.Branch + ":" + gh.args.Path)
	gh.QueryVariables["repository"] = github.String(gh.args.Repository)
}

func (gh *GHGetFile) Run() (int, []byte, error) {
	var err error
	output := Outputs{}
	query := &fileContentsQuery

	err = gh.Client.Query(context.Background(), &query, gh.QueryVariables)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get-file failed with error: error: %w", err)
	}

	output.FileContent = query.Repository.Object.Blob.Text
	marshalled, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get-file failed with error: error: %w", err)
	}
	return step.ExitCodeOK, marshalled, err
}

func main() {
	step.Run(&GHGetFile{})
}

var fileContentsQuery struct {
	Repository struct {
		Object struct {
			Blob struct {
				Text string
			} `graphql:"... on Blob"`
		} `graphql:"object(expression: $filePath)"`
	} `graphql:"repository(name: $repository, owner: $owner)"`
}
