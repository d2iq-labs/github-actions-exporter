package github

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v48/github"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

const githubTokenEnv = "GITHUB_TOKEN"
const scrapeLastEnv = "SCRAPE_LAST_HOURS"
const defaultScrapeLastHours = time.Duration(-48)

// NewClient creates instance of github REST API v3 client
func NewClient(ctx context.Context) (*Client, error) {
	githubToken := os.Getenv(githubTokenEnv)
	if githubToken == "" {
		return nil, errors.New("github token is not provided")
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	httpClient := oauth2.NewClient(ctx, src)

	return &Client{
		client: github.NewClient(httpClient),
	}, nil
}

type WorkflowRuns []*github.WorkflowRun

// GetAllWorkflowRuns fetches all the workflow runs for past X hours
// https://docs.github.com/en/rest/actions/workflow-runs?apiVersion=2022-11-28
func (c *Client) GetAllWorkflowRuns(ctx context.Context, owner, repo string) (WorkflowRuns, error) {
	log.Println("Fetching all workflowruns")
	lastHours, err := getScrapeUntilDuration()
	if err != nil {
		return nil, fmt.Errorf("unable to parse value of past scrape interval hours:%w", err)
	}
	window_start := time.Now().Add(lastHours * time.Hour).Format(time.RFC3339)
	listOpts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 200},
		Created:     ">=" + window_start,
	}
	var runs WorkflowRuns
	for {
		resp, rr, err := c.client.Actions.ListRepositoryWorkflowRuns(context.Background(), owner, repo, listOpts)
		if rl_err, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("ListRepositoryWorkflowRuns ratelimited. until %s", rl_err.Rate.Reset.Time.String())
		} else if err != nil {
			log.Printf("ListRepositoryWorkflowRuns error for repo %s/%s: %s", owner, repo, err.Error())
			return runs, err
		}

		runs = append(runs, resp.WorkflowRuns...)
		if rr.NextPage == 0 {
			break
		}
		listOpts.Page = rr.NextPage
	}

	return runs, nil
}

type Workflows []*github.Workflow

// Workflows gets all workflow datails
// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28
func (c *Client) GetAllWorkflows(ctx context.Context, owner, repo string) (Workflows, error) {
	log.Println("Fetching all workflows")
	listOpts := &github.ListOptions{
		PerPage: 200,
	}
	var flows Workflows
	for {
		resp, rr, err := c.client.Actions.ListWorkflows(context.Background(), owner, repo, listOpts)
		if rl_err, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("ListWorkflows ratelimited. until %s", rl_err.Rate.Reset.Time.String())
		} else if err != nil {
			log.Printf("ListWorkflows error for repo %s/%s: %s", owner, repo, err.Error())
			return flows, err
		}

		flows = append(flows, resp.Workflows...)
		if rr.NextPage == 0 {
			break
		}
		listOpts.Page = rr.NextPage
	}

	return flows, nil
}

type WorkflowsJobs []*github.WorkflowJob

// GetAllWorkflowJobs gets all workflow jobs
// https://docs.github.com/en/rest/actions/workflows?apiVersion=2022-11-28
func (c *Client) GetAllWorkflowJobs(ctx context.Context, owner, repo string) (WorkflowsJobs, error) {
	log.Println("Fetching all workflowjobs")
	workflowRuns, err := c.GetAllWorkflowRuns(ctx, owner, repo)
	if err != nil {
		return nil, err
	}
	var jobs WorkflowsJobs
	for _, wr := range workflowRuns {
		runID := wr.ID
		jobsForRun, err := c.getAllWorkflowJobs(ctx, owner, repo, *runID)
		if err != nil {
			return nil, fmt.Errorf("unable to get all jobs for run: %s : %s", wr.GetName(), wr.GetHTMLURL())
		}
		jobs = append(jobs, jobsForRun...)
	}
	return jobs, nil

}

func (c *Client) getAllWorkflowJobs(ctx context.Context, owner, repo string, runID int64) (WorkflowsJobs, error) {
	log.Println("Fetching all workflow job for run: ", runID)
	listOpts := &github.ListWorkflowJobsOptions{
		Filter:      "all",
		ListOptions: github.ListOptions{PerPage: 200},
	}
	var jobs WorkflowsJobs
	for {
		resp, rr, err := c.client.Actions.ListWorkflowJobs(context.Background(), owner, repo, runID, listOpts)
		if rl_err, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("ListWorkflowJobs ratelimited. until %s", rl_err.Rate.Reset.Time.String())

		} else if err != nil {
			log.Printf("ListWorkflowJobs error for repo %s/%s: %s", owner, repo, err.Error())
			return jobs, err
		}

		jobs = append(jobs, resp.Jobs...)
		if rr.NextPage == 0 {
			break
		}
		listOpts.Page = rr.NextPage
	}

	return jobs, nil
}

func getScrapeUntilDuration() (time.Duration, error) {

	if envLast := os.Getenv(scrapeLastEnv); envLast != "" {
		return time.ParseDuration(envLast)
	}
	return defaultScrapeLastHours, nil
}
