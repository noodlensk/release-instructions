package main

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	jira "github.com/andygrunwald/go-jira"
	bitbucket "github.com/noodlensk/go-bitbucket"
)

type Ticket struct {
	Key     string
	Name    string
	Commits []Commit
}
type Commit struct {
	ID      string
	Author  string
	Message string
	Repo    *Repo
}
type Result struct {
	Tickets         []Ticket
	HomelessCommits []Commit
	AffectedRepos   []Repo
}
type App struct {
	jiraClient      *jira.Client
	bitbucketClient *bitbucket.Client
	repos           []Repo
}

func NewApp(jiraClient *jira.Client, bitbucketClient *bitbucket.Client, repos []Repo) *App {
	return &App{
		jiraClient:      jiraClient,
		bitbucketClient: bitbucketClient,
		repos:           repos,
	}
}

// TODO: fix or remove context propagation
func (a *App) GenerateRI(ctx context.Context, project, fixVersion string) (*Result, error) {
	var (
		tickets         []Ticket
		affectedRepos   []Repo
		homelessCommits []Commit
	)

	jqlQ := fmt.Sprintf("project='%s' AND fixVersion='%s'", project, fixVersion)
	issues, _, err := a.jiraClient.Issue.Search(jqlQ, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search in JIRA")
	}
	for _, issue := range issues {
		ticket := Ticket{
			Name: issue.Fields.Summary,
			Key:  issue.Key,
		}
		tickets = append(tickets, ticket)
	}

	repoCh := make(chan Repo)
	commitCh := make(chan Commit)
	g, gCtx := errgroup.WithContext(ctx)
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			for {
				select {
				case r, ok := <-repoCh:
					if !ok {
						// channel closed
						return nil
					}
					opt := &bitbucket.CommitsOptions{
						Owner:       config.Bitbucket.Owner,
						Repo_slug:   r.Name,
						Exclude:     r.ExcludeBranch,
						Branchortag: r.Branch,
					}
					commitsResponse, err := a.bitbucketClient.Repositories.Commits.GetCommits(opt)
					if err != nil {
						return errors.Wrapf(err, "failed to get bitbucket commits for %s", r.Name)
					}

					commitsResponseMap, _ := commitsResponse.(map[string]interface{})
					commitsResponseValuesList, isList := commitsResponseMap["values"].([]interface{})
					if !isList {
						return errors.Errorf("failed to conver bitbucket response to list")
					}
					for _, v := range commitsResponseValuesList {
						vMap, isMap := v.(map[string]interface{})
						if !isMap {
							return errors.Errorf("failed to conver bitbucket response to map")
						}
						msg := vMap["message"].(string)
						commitCh <- Commit{
							Message: msg,
							Repo:    &r,
						}
					}
				case <-gCtx.Done():
					return gCtx.Err()
				}
			}
		})

	}
	go func() {
		if err := g.Wait(); err != nil {
			log.Error(err)
		}
		close(commitCh)
	}()
	var otherCommits []Commit
	var affectedReposMap = make(map[string]*Repo)
	gRes := errgroup.Group{}
	gRes.Go(func() error {
		for r := range commitCh {
			for i, t := range tickets {
				if strings.Contains(r.Message, t.Key) {
					tickets[i].Commits = append(t.Commits, r)
					affectedReposMap[r.Repo.Name] = r.Repo
				} else {
					otherCommits = append(otherCommits, r)
				}
			}
		}
		return nil
	})
	for _, repo := range a.repos {
		repoCh <- repo
	}
	close(repoCh)
	gRes.Wait()

	for _, c := range otherCommits {
		for r := range affectedReposMap {
			if c.Repo.Name == r {
				homelessCommits = append(homelessCommits, c)
			}
		}
	}
	for _, r := range affectedReposMap {
		affectedRepos = append(affectedRepos, *r)
	}
	return &Result{
		Tickets:         tickets,
		AffectedRepos:   affectedRepos,
		HomelessCommits: homelessCommits,
	}, nil
}
