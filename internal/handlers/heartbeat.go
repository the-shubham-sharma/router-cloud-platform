package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"router-cloud-platform/internal/cache"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/models"
	"router-cloud-platform/internal/utils"
	ws "router-cloud-platform/internal/websocket"
)

type HeartbeatRequest struct {
	CPU       float64 `json:"cpu" binding:"required"`
	Memory    float64 `json:"memory" binding:"required"`
	Latency   float64 `json:"latency" binding:"required"`
	Bandwidth float64 `json:"bandwidth" binding:"required"`
}

func SendHeartbeat(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("id")

	var device models.Device
	if err := database.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Device not found")
		return
	}

	var req HeartbeatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	heartbeat := models.Heartbeat{
		ID:        uuid.New(),
		DeviceID:  device.ID,
		CPU:       req.CPU,
		Memory:    req.Memory,
		Latency:   req.Latency,
		Bandwidth: req.Bandwidth,
	}

	if err := database.DB.Create(&heartbeat).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to save heartbeat")
		return
	}

	database.DB.Model(&device).Update("status", models.StatusOnline)

	cacheKey := "heartbeat:" + deviceID
	data, _ := json.Marshal(heartbeat)
	cache.Client.Set(context.Background(), cacheKey, data, 2*time.Minute)

	ws.BroadcastDeviceUpdate(deviceID, string(models.StatusOnline), req.CPU, req.Memory)

	utils.Success(c, http.StatusOK, "Heartbeat received", heartbeat)
}

func GetLatestHeartbeat(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("id")

	var device models.Device
	if err := database.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Device not found")
		return
	}

	cacheKey := "heartbeat:" + deviceID
	cached, err := cache.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var heartbeat models.Heartbeat
		json.Unmarshal([]byte(cached), &heartbeat)
		utils.Success(c, http.StatusOK, "Latest heartbeat (cache)", heartbeat)
		return
	}

	var heartbeat models.Heartbeat
	if err := database.DB.Where("device_id = ?", deviceID).
		Order("created_at DESC").
		First(&heartbeat).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "No heartbeat found")
		return
	}

	utils.Success(c, http.StatusOK, "Latest heartbeat", heartbeat)
}