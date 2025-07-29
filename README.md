
# 🧭 Distributed Job Scheduler


## 📦 Technologies Used

- Go (Golang)
- SQL Server
- Redis
- Docker 
- go-redis
- godotenv

## 🚀 Features

- Enqueue jobs with priority and retry limits
- Redis queue for task buffering
- SQL Server for persistent job metadata
- Worker pool for concurrent job execution
- Horizontal scaling via container replicas
- Failed job logging and reprocessing


---

## 📐 High-Level Architecture

A horizontally scalable distributed job scheduler in Go, using:
- Native Go concurrency (goroutines, channels)
- Priority queueing
- Retry and Dead Letter Queue (DLQ) support
- MSSQL integration
- Extensible worker pool
- Has Docker Image For  Horizontal scaling (5 workers)
```
+------------------+             +---------------------+
|  Job Submitter   | --------->  |   Jobs Table (DB)   |
|  (API)     |             +---------------------+
       |
       v
+------------------+                                +-----------------------+
| Scheduler Engine |                                |   FailedJobs Table    |
| - Loads jobs     |                                |   (Dead Letter Queue) |
| - Enqueues based |                                +-----------------------+
|   on priority    |                                           ^
       |                                                     |
       v                                                     |
+------------------+            +------------------+         |
| Job Queue System |  <------>  |  Worker Pool      | -------+
| - 3 priority Qs   |           | - N workers       |
| - Retry logic     |           | - Retry & fail    |
+------------------+           +------------------+
```

---
## 🌐 Services Overview

| Service       | Description               | Port     |
|---------------|---------------------------|----------|
| Redis         | Task queue                | 6379     |
| SQL Server    | Job metadata persistence  | 1433     |
| Go App        | Main scheduler service    | 8080     |

---
## 📦 Required Go Packages

Install using:

```bash
go get github.com/go-redis/redis/v9
go get gorm.io/driver/sqlserver
go get github.com/joho/godotenv
go get github.com/google/uuid
```
## 🧩 High-Level Implementation (Responsibilities)

### ✅ Job Submitter
- Accepts job submission (via API)
- Inserts into DB with status = `queued`

### ✅ Scheduler Engine (`main.go`)
- Periodically fetches pending jobs from DB
- Enqueues them based on priority

### ✅ Job Queue (`queue/queue.go`)
- Manages 3 queues: `High`, `Medium`, `Low`
- Handles fair dequeueing

### ✅ Worker Pool (`worker/worker.go`)
- Concurrently processes jobs
- On failure:
  - Retries if available
  - Sends to DLQ if retries exhausted

### ✅ Database Layer (`storage.go`)
- Fetching jobs
- Updating status/retries
- Inserting to DLQ

---

## 📚 Queue Design

### 🗂️ Priority Queues

| Priority  |  
|-----------|
| High      | 
| Medium    | 
| Low       | 

### 🔁 Dequeue Logic (Pseudocode)

```go
function Dequeue():
  if High has job:
    return job
  else if Medium has job:
    return job
  else:
    return job from Low
```

---

## 💣 Dead Letter Queue (DLQ)
- After retries exhausted, jobs are moved to `FailedJobs` table
- Tracked with failure reason

---

## 🔁 Retry Flow

```
[Job Picked]
     |
     v
[Process Job]
     |
     v
[Error?] ---- No ---> [Mark COMPLETED]
     |
    Yes
     |
[Retries Left > 0?] --- No ---> [Send to DLQ + Mark FAILED]
     |
    Yes
     |
[Decrement Retry]
     |
[Requeue Job]
```

---



## 🚀 How to Run With Docker
set your mssql password for mssql image in the docker-compse.yml file 
with docker command 

docker-compose up --build

### 1. Configure MSSQL Connection

## Database Scripts a is in the Folder Called SQL_File
Execute  that File and  then Edit `.env` File:

```
DB_USER=youruser
DB_PASS=yourpass
DB_NAME=yourdb
DB_HOST=localhost
DB_PORT=1433
```


### 2. Run the Project

```bash
go mod tidy
go run cmd/main.go
```

Jobs will be fetched every 2 minutes and processed.

---

## 🧪 Testing API (Postman)
You can test the APIs via Postman or cURL.
Use the Url `http://localhost:8080/`:

- `POST /enqueue`: Add job to DB
- `GET /getAllJobs`: List all jobs
- `GET /failed-jobs`: List DLQ

---
## Sample Data for Post Json Data 
  [{
    "type": "send_email",
    "payload": "{\"to\":\"user1@example.com\",\"subject\":\"Welcome!\"}",
    "priority": "high",
    "max_retries": 3
  },
{
    "type": "export_user_data",
    "payload": "{\"to\":\"user1@example.com\",\"subject\":\"Welcome!\"}",
    "priority": "high",
    "max_retries": 3
  },
]

## 🛠 Future Enhancements

- REST API for job submission
- Real workers (email, PDF, etc.)
- Dashboard UI for DLQ/job tracking
- Docker/K8s deployment
