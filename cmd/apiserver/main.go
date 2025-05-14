package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kylej21/CI-CD-Runner/internal/jobs"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "CI Runner OK")
	})
	http.HandleFunc("/job", func(w http.ResponseWriter, r *http.Request) {
		var job jobs.Job = jobs.Job{
			ID:      "node-test",
			RepoURL: "https://github.com/heroku/node-js-sample.git",
			Branch:  "main",
		}
		job.Start()
		fmt.Fprintf(w, "Job created: %+v\n", job)
		fmt.Fprintf(w, "Job status:  %+v\n", job.Status)
		job.Run()
		fmt.Fprintf(w, "Build status: %+v\n", job.Status)
	})
	log.Println("API server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
