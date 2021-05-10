package main

import (
	"encoding/json"
	"fmt"

	"github.com/andygrunwald/go-jira"
	envconf "github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/atlassian/jira/base"
	"golang.org/x/oauth2"
)

type Args struct {
	base.Args
	IssueSummary     string `env:"ISSUE_SUMMARY,required"`
	JiraProject      string `env:"JIRA_PROJECT,required"`
	JiraCategory     string `env:"ISSUE_CATEGORY"`
	IssueDescription string `env:"ISSUE_DESCRIPTION"`
	IssueType        string `env:"ISSUE_TYPE,required"`
	AssigneeEmail    string `env:"ASSIGNEE_EMAIL,required"`
	ReporterEmail    string `env:"REPORTER_EMAIL"`
}

type jiraCreateIssue struct {
	args  Args
	token *oauth2.Token
}

type output struct {
	Id  string `json:"id"`
	Key string `json:"key"`
	step.Outputs
}

func (j *jiraCreateIssue) Init() error {
	err := envconf.Parse(&j.args)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(j.args.JiraToken), &j.token); err != nil {
		return fmt.Errorf("unmarshal token: %w", err)
	}

	return nil
}

func (j *jiraCreateIssue) Run() (int, []byte, error) {
	oauthClient, url, err := base.GetOauthConnectionDetails(j.token, j.args.SiteName)
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	jiraClient, err := jira.NewClient(oauthClient, url)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("create jira client: %w", err)
	}

	assignee, resp, err := jiraClient.User.Find(j.args.AssigneeEmail)
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("find assignee: %w", err)
	}

	if len(assignee) != 1 {
		return step.ExitCodeFailure, nil, fmt.Errorf("couldn't find assignee, please check that the users exist")
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Assignee: &jira.User{
				AccountID: assignee[0].AccountID,
			},
			Description: j.args.IssueDescription,
			Type: jira.IssueType{
				Name: j.args.IssueType,
			},
			Project: jira.Project{
				Key: j.args.JiraProject,
				ProjectCategory: jira.ProjectCategory{
					Name: j.args.JiraCategory,
				},
			},
			Summary: j.args.IssueSummary,
		},
	}

	if j.args.ReporterEmail != "" {
		reporter, resp, err := jiraClient.User.Find(j.args.ReporterEmail)
		if err != nil {
			return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("find reporter: %w", err)
		}

		if len(reporter) != 1 {
			return step.ExitCodeFailure, nil, fmt.Errorf("couldn't find reporter, please check that the users exist")
		}

		i.Fields.Reporter = &jira.User{
			AccountID: reporter[0].AccountID,
		}
	}

	createdIssue, resp, err := jiraClient.Issue.Create(&i)
	if err != nil {
		return step.ExitCodeFailure, base.ExtractResponse(resp), fmt.Errorf("create issue: %w", err)
	}

	jsonOutput, err := json.Marshal(&output{
		Outputs: step.Outputs{Object: createdIssue},
		Id:      createdIssue.ID,
		Key:     createdIssue.Key,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, err
	}

	return step.ExitCodeOK, jsonOutput, nil
}

func main() {
	step.Run(&jiraCreateIssue{})
}
