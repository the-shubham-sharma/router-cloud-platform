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
)

type CreateDeviceRequest struct {
	Name      string `json:"name" binding:"required"`
	IPAddress string `json:"ip_address" binding:"required"`
	Location  string `json:"location"`
}

type UpdateDeviceRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Location  string `json:"location"`
}

func CreateDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	device := models.Device{
		ID:        uuid.New(),
		UserID:    userID.(uuid.UUID),
		Name:      req.Name,
		IPAddress: req.IPAddress,
		Location:  req.Location,
		Status:    models.StatusUnknown,
	}

	if err := database.DB.Create(&device).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to create device")
		return
	}

	// Invalidate cache
	cache.Client.Del(context.Background(), "devices:"+userID.(uuid.UUID).String())

	utils.Success(c, http.StatusCreated, "Device created", device)
}

func GetDevices(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cacheKey := "devices:" + userID.(uuid.UUID).String()

	// Check cache
	cached, err := cache.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var devices []models.Device
		json.Unmarshal([]byte(cached), &devices)
		utils.Success(c, http.StatusOK, "Devices fetched (cache)", devices)
		return
	}

	var devices []models.Device
	if err := database.DB.Where("user_id = ?", userID).Find(&devices).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch devices")
		return
	}

	// Store in cache for 5 minutes
	data, _ := json.Marshal(devices)
	cache.Client.Set(context.Background(), cacheKey, data, 5*time.Minute)

	utils.Success(c, http.StatusOK, "Devices fetched", devices)
}

func GetDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("id")

	var device models.Device
	if err := database.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Device not found")
		return
	}

	utils.Success(c, http.StatusOK, "Device fetched", device)
}

func UpdateDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("id")

	var device models.Device
	if err := database.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Device not found")
		return
	}

	var req UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Name != "" {
		device.Name = req.Name
	}
	if req.IPAddress != "" {
		device.IPAddress = req.IPAddress
	}
	if req.Location != "" {
		device.Location = req.Location
	}

	if err := database.DB.Save(&device).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to update device")
		return
	}

	// Invalidate cache
	cache.Client.Del(context.Background(), "devices:"+userID.(uuid.UUID).String())

	utils.Success(c, http.StatusOK, "Device updated", device)
}

func DeleteDevice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deviceID := c.Param("id")

	var device models.Device
	if err := database.DB.Where("id = ? AND user_id = ?", deviceID, userID).First(&device).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "Device not found")
		return
	}

	if err := database.DB.Delete(&device).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to delete device")
		return
	}

	// Invalidate cache
	cache.Client.Del(context.Background(), "devices:"+userID.(uuid.UUID).String())

	utils.Success(c, http.StatusOK, "Device deleted", nil)
}