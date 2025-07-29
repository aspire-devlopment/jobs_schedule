package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

// Global DB connection
var db *sql.DB

// Initialize database connection at startup
func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env variables.")
	}

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable",
		user, pass, host, port, name)

	var err error
	db, err = sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("DB ping failed: %v", err)
	}

	log.Println("Connected to MSSQL successfully")
}

// SaveJob inserts a new job record into the Jobs table
func SaveJob(job Job) error {
	query := `
		INSERT INTO Jobs 
		(Id, Type, Payload, Priority, MaxRetries, RetryCount, Status, CreatedAt, UpdatedAt)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, @p8, @p9)
	`

	_, err := db.Exec(query,
		job.Id, job.Type, job.Payload,
		job.Priority, job.MaxRetries, job.RetryCount,
		job.Status, job.CreatedAt, job.UpdatedAt)

	if err != nil {
		log.Printf("Failed to save job ID=%s: %v", job.Id, err)
	}

	return err
}

// UpdateJob updates RetryCount, Status, and UpdatedAt of a job
func UpdateJob(job Job) error {
	job.UpdatedAt = time.Now()

	query := `
		UPDATE Jobs
		SET RetryCount = @p1,
			Status     = @p2,
			UpdatedAt  = @p3
		WHERE Id = @p4
	`

	_, err := db.Exec(query, job.RetryCount, job.Status, job.UpdatedAt, job.Id)

	if err != nil {
		log.Printf("Failed to update job ID=%s: %v", job.Id, err)
	}

	return err
}

// SaveFailedJob inserts a job that has failed permanently into the FailedJobs table
func SaveFailedJob(failed FailedJob) error {
	log.Printf("Inserting Failed Job:\n"+
		"ID=%s, OriginalJobId=%s, Type=%s, Priority=%s, Reason=%s, FailedAt=%s\nPayload: %s\n",
		failed.Id, failed.OriginalJobId, failed.Type,
		failed.Priority, failed.Reason, failed.FailedAt.Format(time.RFC3339), failed.Payload)

	query := `
		INSERT INTO FailedJobs 
		(Id, OriginalJobId, Type, Payload, Priority, Reason, FailedAt)
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7)
	`

	_, err := db.Exec(query,
		failed.Id, failed.OriginalJobId, failed.Type,
		failed.Payload, failed.Priority, failed.Reason, failed.FailedAt)

	if err != nil {
		log.Printf("Error inserting failed job ID=%s: %v", failed.Id, err)
	}

	return err
}

// GetAllFailedJobs returns a list of all failed jobs
func GetAllFailedJobs() ([]FailedJob, error) {
	query := `
		SELECT 
			CONVERT(VARCHAR(36), Id) AS Id, OriginalJobId, Type, Payload, Priority, Reason, FailedAt
		FROM FailedJobs`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var failedJobs []FailedJob
	for rows.Next() {
		var f FailedJob
		if err := rows.Scan(&f.Id, &f.OriginalJobId, &f.Type, &f.Payload, &f.Priority, &f.Reason, &f.FailedAt); err != nil {
			log.Println("Failed to scan failed job:", err)
			continue
		}
		failedJobs = append(failedJobs, f)
	}
	return failedJobs, nil
}

// getAllJobs returns a list of all  jobs
func getAllJobs() ([]Job, error) {
	query := `Select
		CONVERT(VARCHAR(36), Id) AS Id,
			Type, Payload, Status, Priority,
			RetryCount, MaxRetries
		FROM Jobs`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.Id, &job.Type, &job.Payload, &job.Status, &job.Priority, &job.RetryCount, &job.MaxRetries); err != nil {
			log.Println("Failed to scan pending job:", err)
			continue
		}
		log.Printf("Loaded pending job: %+v", job)
		jobs = append(jobs, job)
	}

	return jobs, nil
}

// GetPendingJobsFromDB fetches jobs that are still pending (queued or retrying)
func GetPendingJobsFromDB() ([]Job, error) {
	query := `
		SELECT 
			CONVERT(VARCHAR(36), Id) AS Id,
			Type, Payload, Status, Priority,
			RetryCount, MaxRetries
		FROM Jobs
		WHERE Status IN ('queued', 'retrying')
		ORDER BY 
			CASE Priority
				WHEN 1 THEN 1  -- high
				WHEN 2 THEN 2  -- medium
				WHEN 3 THEN 3  -- low
				ELSE 4
			END,
			CreatedAt ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		if err := rows.Scan(&job.Id, &job.Type, &job.Payload, &job.Status, &job.Priority, &job.RetryCount, &job.MaxRetries); err != nil {
			log.Println("Failed to scan pending job:", err)
			continue
		}
		log.Printf("Loaded pending job: %+v", job)
		jobs = append(jobs, job)
	}

	return jobs, nil
}
