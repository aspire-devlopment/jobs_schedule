package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// to save the metada into jobs table  can insert multiple json at one time
func EnqueueJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var jobs []Job
	if err := json.NewDecoder(r.Body).Decode(&jobs); err != nil {
		http.Error(w, "Invalid JSON format. Expected an array of jobs.", http.StatusBadRequest)
		return
	}

	var successJobs []Job
	var failedJobs []map[string]string

	for _, job := range jobs {
		if err := validateJob(job); err != nil {
			failedJobs = append(failedJobs, map[string]string{
				"type":  job.Type,
				"error": err.Error(),
			})
			continue
		}
		job.Id = uuid.New().String()
		job.Status = "queued"
		job.RetryCount = 0
		job.CreatedAt = time.Now()
		job.UpdatedAt = time.Now()

		if err := SaveJob(job); err != nil {
			failedJobs = append(failedJobs, map[string]string{
				"type":  job.Type,
				"error": "Failed to save job",
			})
			continue
		}

		if err := Enqueue(job); err != nil {
			failedJobs = append(failedJobs, map[string]string{
				"type":  job.Type,
				"error": "Failed to enqueue job",
			})
			continue
		}

		successJobs = append(successJobs, job)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": successJobs,
		"failed":  failedJobs,
	})
}

func validateJob(job Job) error {

	// Validate Payload is valid JSON
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(job.Payload), &payloadMap); err != nil {
		return fmt.Errorf("invalid payload JSON")
	}

	// Type-specific payload validation
	if job.Type == "send_email" {
		if _, ok := payloadMap["to"]; !ok {
			return fmt.Errorf("missing 'to' field in payload")
		}
		if _, ok := payloadMap["subject"]; !ok {
			return fmt.Errorf("missing 'subject' field in payload")
		}
	}

	// Validate Priority
	switch job.Priority {
	case "low", "medium", "high":
		// valid
	default:
		return fmt.Errorf("invalid priority: must be 'low', 'medium', or 'high'")
	}

	// Validate MaxRetries
	if job.MaxRetries < 0 {
		return fmt.Errorf("max_retries must be zero or more")
	}

	return nil
}

func GetFailedJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	failedJobs, err := GetAllFailedJobs()
	if err != nil {
		http.Error(w, "Failed to get failed jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(failedJobs)
}
func GetAllJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	allJobs, err := getAllJobs()
	if err != nil {
		http.Error(w, "Failed to get failed jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allJobs)
}
