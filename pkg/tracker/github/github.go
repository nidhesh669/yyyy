package github

import (
	"context"

	"strconv"

	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/tracker/config"
	"golang.org/x/oauth2"
)

const (
	name = "github"
)

// code snippet https://github.com/while-loop/todo/blob/cc6b554cccfd3598f6b6342d69c78abcbc5d0128/pkg/app.go#L17-L25
// footer  ###### This issue was generated by [todo](https://github.com/while-loop/todo) on behalf of %s.
type Tracker struct {
	conf     *config.GithubConfig
	ghClient *github.Client
}

func NewTracker(cfg *config.GithubConfig) *Tracker {
	// todo create ghclient with accesstoken from github app
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.AccessToken})
	oauthClient := oauth2.NewClient(context.Background(), ts)
	return &Tracker{
		conf:     cfg,
		ghClient: github.NewClient(oauthClient),
	}
}

func (t *Tracker) GetIssues(ctx context.Context, owner, repo string) ([]*issue.Issue, error) {
	gIss, _, err := t.ghClient.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get issues for %s/%s", owner, repo)
	}

	// todo support pagination
	iss := make([]*issue.Issue, 0)

	for _, is := range gIss {
		iss = append(iss, ghIssue2todoIssue(owner, repo, is))
	}

	return iss, nil
}

func (t *Tracker) CreateIssue(ctx context.Context, issue *issue.Issue) (*issue.Issue, error) {
	is, _, err := t.ghClient.Issues.Create(ctx, issue.Owner, issue.Repo, &github.IssueRequest{
		Title:    &issue.Title,
		Body:     &issue.Description,
		Labels:   &issue.Labels,
		Assignee: &issue.Assignee,
		State:    pString("open"),
	})

	if err != nil {
		return nil, errors.Wrapf(err, "failed to create issue %s/%s", issue.Owner, issue.Repo)
	}

	return ghIssue2todoIssue(issue.Owner, issue.Repo, is), nil
}

func (t *Tracker) DeleteIssue(ctx context.Context, issue *issue.Issue) error {
	iID, _ := strconv.Atoi(issue.ID)
	_, resp, err := t.ghClient.Issues.Edit(ctx, issue.Owner, issue.Repo, iID, &github.IssueRequest{
		State: pString("closed"),
	})

	if err != nil {
		return errors.Wrapf(err, "failed to close issue %s/%s/%d", issue.Owner, issue.Repo, iID)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to close issue %s/%s/%d.. http status: %d", issue.Owner, issue.Repo, iID, resp.StatusCode)
	}

	return nil
}

func (t *Tracker) Name() string {
	return name
}

func parseLabels(gLs []github.Label) []string {
	ls := []string{}
	for _, gL := range gLs {
		ls = append(ls, gL.GetName())
	}
	return ls
}

func pString(s string) *string {
	return &s
}

func ghIssue2todoIssue(owner, repo string, gIs *github.Issue) *issue.Issue {
	return &issue.Issue{
		ID:          strconv.Itoa(gIs.GetID()),
		Title:       gIs.GetTitle(),
		Description: gIs.GetBody(),
		Assignee:    gIs.GetAssignee().GetName(),
		Author:      gIs.GetUser().GetName(),
		Mentions:    []string{},
		Labels:      parseLabels(gIs.Labels),
		File:        "",
		Line:        0,
		Owner:       owner,
		Repo:        repo,
	}
}
