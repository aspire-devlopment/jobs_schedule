package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// Shared context and Redis client (for all workers)
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // You can parametrize this for flexibility
})

// Dequeue pulls up to `batchSize` jobs, prioritizing high → medium → low queues
func Dequeue(batchSize int) []Job {
	var jobs []Job
	priorities := []string{"high", "medium", "low"}

	for _, p := range priorities {
		queueKey := "queue:jobs:" + p

		// Fill only the remaining slot of batchSize
		for i := 0; i < batchSize-len(jobs); i++ {
			result, err := rdb.RPop(ctx, queueKey).Result()
			if err != nil || result == "" {
				if err != redis.Nil && err != nil {
					log.Printf("Redis error while popping from %s: %v\n", queueKey, err)
				}
				break // stop trying this queue
			}

			var job Job
			if err := json.Unmarshal([]byte(result), &job); err != nil {
				log.Printf("Failed to unmarshal job from Redis (%s): %v\n", queueKey, err)
				continue
			}

			log.Printf("Dequeued job: ID=%s, Type=%s, Priority=%s\n", job.Id, job.Type, job.Priority)
			jobs = append(jobs, job)
		}

		// Exit early if we've already got a full batch
		if len(jobs) >= batchSize {
			break
		}
	}

	if len(jobs) == 0 {
		log.Println("No jobs dequeued this round.")
	}
	return jobs
}

// Enqueue adds a job to Redis with deduplication logic using a Redis Set
func Enqueue(job Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	if err != nil {
		log.Printf("Redis error during dedup check for job ID=%s: %v\n", job.Id, err)
		return err
	}

	queueKey := "queue:jobs:" + job.Priority

	// Use pipeline to perform both LPUSH and SADD atomically
	pipe := rdb.TxPipeline()
	pipe.LPush(ctx, queueKey, data)
	pipe.SAdd(ctx, "dedup:jobs", job.Id)
	_, err = pipe.Exec(ctx)

	if err != nil {
		log.Printf("Failed to enqueue job ID=%s: %v\n", job.Id, err)
		return err
	}
	log.Printf("Enqueued job ID=%s to queue: %s\n", job.Id, queueKey)
	return nil
}

// LoadPendingJobsFromDBToRedis pulls jobs from DB and enqueues them into Redis
func LoadPendingJobsFromDBToRedis() {
	log.Println("Syncing pending jobs from DB to Redis...")
	jobs, err := GetPendingJobsFromDB()
	if err != nil {
		LogError("Failed to fetch pending jobs from DB: " + err.Error())
		return
	}

	log.Printf("Fetched %d pending jobs from DB\n", len(jobs))
	for _, job := range jobs {
		log.Printf("Syncing job ID=%s (Type=%s, Priority=%s)\n", job.Id, job.Type, job.Priority)
		if err := Enqueue(job); err != nil {
			LogError("Failed to enqueue job from DB: " + err.Error())
		}
	}
}

// SendToDLQ moves failed jobs into a Redis DLQ and logs them in the database
func SendToDLQ(job Job, reason string) error {
	job.Status = "failed"

	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	// Push to Redis DLQ
	if err := rdb.LPush(ctx, "queue:dlq", data).Err(); err != nil {
		log.Printf("Failed to push job ID=%s to DLQ: %v\n", job.Id, err)
		return err
	}

	// save failure dlq in database
	failed := FailedJob{
		Id:            uuid.New().String(),
		OriginalJobId: job.Id,
		Type:          job.Type,
		Payload:       job.Payload,
		Priority:      job.Priority,
		Reason:        reason,
		FailedAt:      time.Now(),
	}
	return SaveFailedJob(failed)
}
