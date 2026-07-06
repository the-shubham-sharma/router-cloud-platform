package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"router-cloud-platform/internal/cache"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/models"
)

type HeartbeatJob struct {
	DeviceID  uuid.UUID
	CPU       float64
	Memory    float64
	Latency   float64
	Bandwidth float64
}

type HeartbeatWorkerPool struct {
	jobs    chan HeartbeatJob
	workers int
}

var Pool *HeartbeatWorkerPool

func NewHeartbeatWorkerPool(workers int, queueSize int) *HeartbeatWorkerPool {
	return &HeartbeatWorkerPool{
		jobs:    make(chan HeartbeatJob, queueSize),
		workers: workers,
	}
}

func (p *HeartbeatWorkerPool) Start() {
	for i := 0; i < p.workers; i++ {
		go p.runWorker(i + 1)
	}
	log.Printf("Heartbeat worker pool started with %d workers", p.workers)
}

func (p *HeartbeatWorkerPool) Submit(job HeartbeatJob) bool {
	select {
	case p.jobs <- job:
		return true
	default:
		log.Println("Worker pool queue full")
		return false
	}
}

func (p *HeartbeatWorkerPool) runWorker(id int) {
	for job := range p.jobs {
		p.processHeartbeat(id, job)
	}
}

func (p *HeartbeatWorkerPool) processHeartbeat(workerID int, job HeartbeatJob) {
	heartbeat := models.Heartbeat{
		ID:        uuid.New(),
		DeviceID:  job.DeviceID,
		CPU:       job.CPU,
		Memory:    job.Memory,
		Latency:   job.Latency,
		Bandwidth: job.Bandwidth,
	}
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		err = database.DB.Create(&heartbeat).Error
		if err == nil {
			break
		}
		time.Sleep(time.Duration(attempt*100) * time.Millisecond)
	}
	if err != nil {
		log.Printf("[Worker %d] Failed after 3 attempts", workerID)
		return
	}
	database.DB.Model(&models.Device{}).Where("id = ?", job.DeviceID).Update("status", models.StatusOnline)
	cacheKey := "heartbeat:" + job.DeviceID.String()
	data, _ := json.Marshal(heartbeat)
	cache.Client.Set(context.Background(), cacheKey, data, 2*time.Minute)
	log.Printf("[Worker %d] Heartbeat saved successfully", workerID)
}