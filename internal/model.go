package internal

import "time"

//models for the tables
type Job struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Payload    string    `json:"payload"`
	Priority   string    `json:"priority"` // "high", "medium", "low"
	MaxRetries int       `json:"max_retries"`
	RetryCount int       `json:"retry_count"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FailedJob struct {
	Id            string    `json:"id"`
	OriginalJobId string    `json:"original_job_id"`
	Type          string    `json:"type"`
	Payload       string    `json:"payload"`
	Priority      string    `json:"priority"`
	Reason        string    `json:"reason"`
	FailedAt      time.Time `json:"failed_at"`
}
