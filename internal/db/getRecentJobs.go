package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/kylej21/CI-CD-Runner/internal/jobs"
)

func GetRecentJobs(conn *pgx.Conn, limit int) ([]jobs.Job, error) {
	// simple query to get job info in reverse order of creation
	rows, err := conn.Query(context.Background(), `
		SELECT id, repo_url, branch, status, created_at, finished_at
		FROM jobs
		ORDER BY created_at DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []jobs.Job

	// built in list traversal
	for rows.Next() {
		var j jobs.Job
		err := rows.Scan(&j.ID, &j.RepoURL, &j.Branch, &j.Status, &j.CreatedAt, &j.FinishedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, j)
	}

	return result, nil
}
