package main

import (
	"context"
	"fmt"
	"log"

	"github.com/d2iq-labs/github-actions-exporter/pkg/github"
)

func main() {
	ctx := context.Background()
	ghClient, err := github.NewClient(ctx)
	if err != nil {
		log.Fatal("unable to initialize github client:", err)
	}
	workflows, err := ghClient.GetAllWorkflows(ctx, "mesosphere", "konvoy2")
	if err != nil {
		log.Fatal(err)
	}
	for _, w := range workflows {
		fmt.Printf("workflow id: %s\n", *w.Name)
	}

	jobs, err := ghClient.GetAllWorkflowJobs(ctx, "mesosphere", "konvoy2")
	if err != nil {
		log.Fatal(err)
	}
	for _, j := range jobs {
		fmt.Printf("workflow job: %s\n", *j.Name)
	}

}
