package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/models"
	"router-cloud-platform/internal/utils"
)

func AdminGetAllDevices(c *gin.Context) {
	var devices []models.Device
	if err := database.DB.Find(&devices).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch devices")
		return
	}
	utils.Success(c, http.StatusOK, "All devices", devices)
}

func AdminGetAllUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	utils.Success(c, http.StatusOK, "All users", users)
}

func AdminPromoteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "User not found")
		return
	}

	if err := database.DB.Model(&user).Update("role", models.RoleAdmin).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to promote user")
		return
	}

	utils.Success(c, http.StatusOK, "User promoted to admin", gin.H{
		"id":    user.ID,
		"email": user.Email,
		"role":  models.RoleAdmin,
	})
}