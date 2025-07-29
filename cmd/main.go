package main

import (
	"distributed_job_scheduler/internal"
	"log"
	"net/http"
	"time"
)

func main() {
	// Initialize the error logger to write logs to the specified file

	internal.InitLogger("logs/error.log")

	// Start periodically loads pending jobs from the database
	// and enqueues them into Redis every 2 minutes if pending queue present in db.
	go func() {
		for {
			internal.LoadPendingJobsFromDBToRedis()
			time.Sleep(2 * time.Minute)
		}
	}()

	// Start a separate function that continuously runs the job worker
	// which dequeues jobs from Redis and processes them.
	go func() {
		for {
			internal.StartWorker() // processes all available jobs in batch
		}
	}()

	//api request to post the metada and to fetch teh failed jobs
	http.HandleFunc("/enqueue", internal.EnqueueJobHandler)
	http.HandleFunc("/failed-jobs", internal.GetFailedJobsHandler)
	http.HandleFunc("/getAllJobs", internal.GetAllJobsHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
