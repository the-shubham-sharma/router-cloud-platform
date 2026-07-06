package alert

import (
	"log"
	"time"

	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/models"
	ws "router-cloud-platform/internal/websocket"
)

func StartOfflineDetector() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		log.Println("Offline device detector started")
		for range ticker.C {
			detectOfflineDevices()
		}
	}()
}

func detectOfflineDevices() {
	twoMinutesAgo := time.Now().Add(-2 * time.Minute)
	var devices []models.Device
	if err := database.DB.Where("status = ?", models.StatusOnline).Find(&devices).Error; err != nil {
		log.Printf("[Alert] Failed to fetch devices: %v", err)
		return
	}
	for _, device := range devices {
		var lastHeartbeat models.Heartbeat
		err := database.DB.Where("device_id = ?", device.ID).
			Order("created_at DESC").
			First(&lastHeartbeat).Error
		if err != nil {
			markOffline(device)
			continue
		}
		if lastHeartbeat.CreatedAt.Before(twoMinutesAgo) {
			markOffline(device)
		}
	}
}

func markOffline(device models.Device) {
	database.DB.Model(&device).Update("status", models.StatusOffline)
	log.Printf("[Alert] Device %s marked OFFLINE", device.Name)
	ws.BroadcastDeviceUpdate(device.ID.String(), string(models.StatusOffline), 0, 0)
}