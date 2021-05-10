package base

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	github "github.com/shurcooL/githubv4"
	"github.com/stackpulse/steps-sdk-go/env"
	"golang.org/x/oauth2"
)

type Args struct {
	Token   string `env:"TOKEN,required" envDefault:""`
	Owner   string `env:"OWNER,required" envDefault:""`
	Privacy string `env:"TYPE" envDefault:""`
}

func (a Args) BaseArgs() Args {
	return a
}

type BaseArgs interface {
	BaseArgs() Args
}

type GithubClient struct {
	Args
	Client         *github.Client
	QueryVariables map[string]interface{}
}

func NewGithubClient(args BaseArgs) (*GithubClient, error) {
	if err := env.Parse(args); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	GHClient := &GithubClient{
		Args:           args.BaseArgs(),
		QueryVariables: map[string]interface{}{},
	}

	if GHClient.Token != "" {
		src := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: GHClient.Token},
		)
		httpClient := oauth2.NewClient(context.Background(), src)
		GHClient.Client = github.NewClient(httpClient)
	} else {
		GHClient.Client = github.NewClient(&http.Client{})
	}

	GHClient.addBaseQueryVariables()

	return GHClient, nil
}

func (GHClient *GithubClient) addBaseQueryVariables() {
	if GHClient.Owner != "" {
		GHClient.QueryVariables["owner"] = github.String(GHClient.Owner)
	}

	if GHClient.Privacy != "" {
		upper := github.RepositoryPrivacy(strings.ToUpper(GHClient.Privacy))
		GHClient.QueryVariables["privacy"] = github.RepositoryPrivacy(upper)
	}
}
