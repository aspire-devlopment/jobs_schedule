package internal

import (
	"fmt"
	"log"
	"time"
)

const (
	BatchSize     = 10            // Number of jobs to fetch per batch from Redis
	WorkerSleep   = time.Minute   // Interval between each worker cycle
	PriorityQueue = "queue:jobs:" // Base Redis queue key prefix for job priorities
)

// StartWorker runs an infinite loop that triggers job processing every WorkerSleep interval.
// we can set it on the minute or hourly basisi from the main.go
// It uses a ticker to schedule periodic job batch processing from Redis.
func StartWorker() {
	ticker := time.NewTicker(WorkerSleep) // Create a ticker that ticks every WorkerSleep duration
	defer ticker.Stop()

	for {
		log.Println("Worker cycle started")
		processAllJobsInBatches() // Process all available jobs in batches

		<-ticker.C // Wait for the next tick before starting next cycle
	}
}

// processAllJobsInBatches continuously fetches batches of jobs from Redis .

func processAllJobsInBatches() {
	for {
		batch := Dequeue(BatchSize) // Fetch up to BatchSize jobs from Redis queue

		if len(batch) == 0 {
			log.Println("No jobs found in Redis. Sleeping before next batch...")
			time.Sleep(3 * time.Second) // Sleep briefly if no jobs found to avoid tight looping
			continue
		}

		log.Printf("Dequeued %d jobs from Redis\n", len(batch))

		// Process each job
		for i := 0; i < len(batch); i++ {
			job := &batch[i]

			log.Printf("Processing job: %s, type: %s\n", job.Id, job.Type)

			// Call the actual job processing function
			err := processJob(job)
			if err != nil {
				// Log error with custom error logger
				LogError("Job " + job.Id + " failed: " + err.Error())

				job.RetryCount++ // Increment retry count

				if job.RetryCount >= job.MaxRetries {
					// If max retries exceeded, mark job as failed and send to dead-letter queue
					job.Status = "failed"
					SendToDLQ(*job, err.Error())
					log.Printf("Job %s moved to DLQ after %d retries\n", job.Id, job.RetryCount)
				} else {
					// Otherwise, mark job as retrying and re-enqueue it for another attempt
					job.Status = "retrying"
					Enqueue(*job) // Re-enqueue failed job
					log.Printf("Job %s re-enqueued for retry %d\n", job.Id, job.RetryCount)
				}
			} else {
				// If job processed successfully, update status to completed
				job.Status = "completed"
				log.Printf("Job %s processed successfully.\n", job.Id)
			}

			// Persist job status and retry count updates to the database
			UpdateJob(*job)
		}
	}
}

// processJob executes the logic for each job based on its Type field.

func processJob(job *Job) error {
	switch job.Type {
	case "send_email":
		log.Println("Sending email with payload:", job.Payload)

		// Simulate failure for demonstration/testing purposes
		return fmt.Errorf("simulated failure for email job")

	case "export_user_data":
		log.Println("Exporting user data with payload:", job.Payload)
		return nil // Indicate success

	case "process_payment":
		log.Println("Processing payment with payload:", job.Payload)
		return nil // Indicate success
	case "Notification":
		log.Println("Sending Notification with payload:", job.Payload)
		return nil // Indicate success
	case "ReportGeneration":
		log.Println("Generating Report with payload:", job.Payload)
		return nil // Indicate success
	default:
		// Unknown job types return an error to trigger failure handling
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}
