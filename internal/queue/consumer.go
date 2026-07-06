package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"router-cloud-platform/internal/cache"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/metrics"
	"router-cloud-platform/internal/models"
)

type HeartbeatMessage struct {
	DeviceID  string  `json:"device_id"`
	CPU       float64 `json:"cpu"`
	Memory    float64 `json:"memory"`
	Latency   float64 `json:"latency"`
	Bandwidth float64 `json:"bandwidth"`
}

func StartConsumer() {
	msgs, err := Channel.Consume(HeartbeatQueue, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}
	log.Println("RabbitMQ consumer started")
	go func() {
		for msg := range msgs {
			processMessage(msg.Body)
			msg.Ack(false)
		}
	}()
}

func processMessage(body []byte) {
	var hb HeartbeatMessage
	if err := json.Unmarshal(body, &hb); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}
	deviceID, err := uuid.Parse(hb.DeviceID)
	if err != nil {
		log.Printf("Invalid device ID: %v", err)
		return
	}
	heartbeat := models.Heartbeat{
		ID:        uuid.New(),
		DeviceID:  deviceID,
		CPU:       hb.CPU,
		Memory:    hb.Memory,
		Latency:   hb.Latency,
		Bandwidth: hb.Bandwidth,
	}
	var dbErr error
	for attempt := 1; attempt <= 3; attempt++ {
		dbErr = database.DB.Create(&heartbeat).Error
		if dbErr == nil {
			break
		}
		log.Printf("[Consumer] DB attempt %d failed: %v", attempt, dbErr)
		time.Sleep(time.Duration(attempt*100) * time.Millisecond)
	}
	if dbErr != nil {
		log.Printf("[Consumer] Failed after 3 attempts")
		return
	}
	database.DB.Model(&models.Device{}).Where("id = ?", deviceID).Update("status", models.StatusOnline)
	cacheKey := "heartbeat:" + hb.DeviceID
	data, _ := json.Marshal(heartbeat)
	cache.Client.Set(context.Background(), cacheKey, data, 2*time.Minute)
	metrics.HeartbeatsTotal.Inc()
	log.Printf("[Consumer] Heartbeat processed for device %s", hb.DeviceID)
}