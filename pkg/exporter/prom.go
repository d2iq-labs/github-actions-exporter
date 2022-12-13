package exporter

import (
	"github.com/google/go-github/v48/github"
	"github.com/prometheus/client_golang/prometheus"
)

func ExportGHAJobs(jobs []*github.WorkflowJob, gv *prometheus.GaugeVec) {
	for _, job := range jobs {
		// "success", = 2
		//       "failure" = 0,
		//       "neutral" = 1,
		//       "cancelled" = 1,
		//       "skipped" = 1,
		//       "timed_out" = 0,
		//       "action_required" = 1,
		//       null = 1
		var metric float64
		switch job.GetConclusion() {
		case "success":
			metric = 2
		case "failure", "timed_out":
			metric = 0
		default:
			metric = 1

		}
		//"org", "project", "job_name", "ref", "status"
		gv.WithLabelValues("mesosphere", "konvoy2", job.GetName(), job.GetHeadSHA(), job.GetConclusion()).Set(metric)

	}
}
