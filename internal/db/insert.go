package db

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kylej21/CI-CD-Runner/internal/jobs"
)

// insert job if it doesn't exist in the local postgres DB
func InsertJob(conn *pgx.Conn, j *jobs.Job, startedAt time.Time) error {
	finishedAt := time.Now()
	duration := finishedAt.Sub(startedAt).Milliseconds()

	var errorCount, warningCount int
	for _, file := range j.LintResults {
		errorCount += file.ErrorCount
		warningCount += file.WarningCount
	}

	buildLog, _ := json.Marshal(j.LintResults)
	_, err := conn.Exec(context.Background(), `
		INSERT INTO jobs (
			id, repo_url, branch, status,
			created_at, finished_at, duration_ms,
			lint_errors, lint_warnings, build_log
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`, j.ID, j.RepoURL, j.Branch, j.Status, startedAt, finishedAt, duration, errorCount, warningCount, string(buildLog),
	)
	return err
}
