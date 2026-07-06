package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/models"
	"router-cloud-platform/internal/utils"
)

func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		var user models.User
		if err := database.DB.First(&user, "id = ?", userID).Error; err != nil {
			utils.Error(c, http.StatusUnauthorized, "User not found")
			c.Abort()
			return
		}
		if user.Role != role {
			utils.Error(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}
		c.Set("user_role", user.Role)
		c.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RequireRole(models.RoleAdmin)
}