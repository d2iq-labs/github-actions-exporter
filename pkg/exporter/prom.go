package exporter

import (
	"log"

	"github.com/google/go-github/v48/github"
	"github.com/prometheus/client_golang/prometheus"
)

func ExportGHAJobs(jobs []*github.WorkflowJob, gv *prometheus.GaugeVec) {
	for _, job := range jobs {
		name := *job.Name
		// "success", = 2
		//       "failure" = 0,
		//       "neutral" = 1,
		//       "cancelled" = 1,
		//       "skipped" = 1,
		//       "timed_out" = 0,
		//       "action_required" = 1,
		//       null = 1
		var metric float64
		if job.Conclusion == nil {
			log.Printf("unable to collect job %s. job's conclusion status is nil", name)
			continue
		}

		switch *job.Conclusion {
		case "success":
			metric = 2
		case "failure", "timed_out":
			metric = 0
		default:
			metric = 1

		}

		gv.WithLabelValues("mesosphere", "konvoy2", name).Set(metric)

	}
}
