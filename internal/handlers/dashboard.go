package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/models"
	"router-cloud-platform/internal/utils"
)

func GetDashboardSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// Total devices
	var totalDevices int64
	database.DB.Model(&models.Device{}).Where("user_id = ?", userID).Count(&totalDevices)

	// Online devices
	var onlineDevices int64
	database.DB.Model(&models.Device{}).Where("user_id = ? AND status = ?", userID, models.StatusOnline).Count(&onlineDevices)

	// Offline devices
	var offlineDevices int64
	database.DB.Model(&models.Device{}).Where("user_id = ? AND status = ?", userID, models.StatusOffline).Count(&offlineDevices)

	// Unknown devices
	var unknownDevices int64
	database.DB.Model(&models.Device{}).Where("user_id = ? AND status = ?", userID, models.StatusUnknown).Count(&unknownDevices)

	// Total heartbeats received
	var totalHeartbeats int64
	database.DB.Model(&models.Heartbeat{}).
		Joins("JOIN devices ON devices.id = heartbeats.device_id").
		Where("devices.user_id = ?", userID).
		Count(&totalHeartbeats)

	// Average CPU across all devices
	type AvgResult struct {
		AvgCPU       float64
		AvgMemory    float64
		AvgLatency   float64
		AvgBandwidth float64
	}

	var avg AvgResult
	database.DB.Model(&models.Heartbeat{}).
		Select("AVG(cpu) as avg_cpu, AVG(memory) as avg_memory, AVG(latency) as avg_latency, AVG(bandwidth) as avg_bandwidth").
		Joins("JOIN devices ON devices.id = heartbeats.device_id").
		Where("devices.user_id = ?", userID).
		Scan(&avg)

	utils.Success(c, http.StatusOK, "Dashboard summary", gin.H{
		"devices": gin.H{
			"total":   totalDevices,
			"online":  onlineDevices,
			"offline": offlineDevices,
			"unknown": unknownDevices,
		},
		"heartbeats": gin.H{
			"total": totalHeartbeats,
		},
		"averages": gin.H{
			"cpu":       avg.AvgCPU,
			"memory":    avg.AvgMemory,
			"latency":   avg.AvgLatency,
			"bandwidth": avg.AvgBandwidth,
		},
	})
}