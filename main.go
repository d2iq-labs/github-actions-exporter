package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/d2iq-labs/github-actions-exporter/pkg/exporter"
	"github.com/d2iq-labs/github-actions-exporter/pkg/github"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	ghClient, err := github.NewClient(ctx)
	if err != nil {
		log.Fatal("unable to initialize github client:", err)
	}
	recordMetrics(ctx, ghClient)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9091", nil)
}

// example from prometheous go guide https://prometheus.io/docs/guides/go-application/
var (
	jobsProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gha_jobs_status",
		Help: "current status of github actions job",
	}, []string{"org", "project", "job_name", "ref", "status"})
)

func recordMetrics(ctx context.Context, ghClient *github.Client) {
	go func() {
		for {
			jobs, err := ghClient.GetAllWorkflowJobs(ctx, "mesosphere", "konvoy2")
			if err != nil {
				log.Println("error fetching github actions job data:", err, "skip collecting data for jobs...")
				continue
			}
			exporter.ExportGHAJobs(jobs, jobsProcessed)
			time.Sleep(300 * time.Second)
		}
	}()
}
