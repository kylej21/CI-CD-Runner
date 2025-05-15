package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/kylej21/CI-CD-Runner/internal/db"
	"github.com/kylej21/CI-CD-Runner/internal/jobs"
)

func main() {

	/*** DB SETUP ----------------------------------- ***/
	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	err = db.SetupDatabase(conn)
	if err != nil {
		log.Fatalf("Failed to setup database schema: %v", err)
	}

	defer conn.Close(context.Background())

	/*** API ENDPOINTS ----------------------------------- ***/

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "CI Runner OK")
	})

	// POST /job
	//
	// Accepts a JSON Job payload, then clones, installs, lints, and builds the specified repo.
	// Returns no response body.
	http.HandleFunc("/job", withCORS(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create an object of type Job
		var job jobs.Job = jobs.Job{
			ID:      "test1",
			RepoURL: "https://github.com/Nico-kun123/starter-pack",
			Branch:  "main",
		}
		fmt.Fprintf(w, "Build started\n")

		// Run the job.
		job.Run()

		// cache results
		err := db.InsertJob(conn, &job, start)
		if err != nil {
			fmt.Fprintf(w, "Failed to save job: %v\n", err)
		} else {
			fmt.Fprintf(w, "Build saved successfully\n")
		}
		fmt.Fprintf(w, "Build status: %+v\n", job.Status)
	}))

	// GET /api/jobs
	//
	// Returns a list of the most recent jobs.
	// Returns a JSON array of Job objects.
	http.HandleFunc("/api/jobs", withCORS(func(w http.ResponseWriter, r *http.Request) {
		jobs, err := db.GetRecentJobs(conn, 20)
		if err != nil {
			http.Error(w, "Failed to fetch jobs", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jobs)
	}))

	log.Println("API server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// defining NO-CORS since local project
func withCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		handler(w, r)
	}
}
