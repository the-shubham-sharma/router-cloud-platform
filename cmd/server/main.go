package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"router-cloud-platform/internal/cache"
	"router-cloud-platform/internal/config"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/handlers"
	"router-cloud-platform/internal/middleware"
	ws "router-cloud-platform/internal/websocket"
)

func main() {
	config.Load()
	database.Connect()
	database.AutoMigrate()
	cache.Connect()

	// Start WebSocket hub
	go ws.GlobalHub.Run()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.GET("/profile", middleware.AuthRequired(), handlers.GetProfile)
	}

	devices := r.Group("/devices", middleware.AuthRequired())
	{
		devices.POST("", handlers.CreateDevice)
		devices.GET("", handlers.GetDevices)
		devices.GET("/:id", handlers.GetDevice)
		devices.PUT("/:id", handlers.UpdateDevice)
		devices.DELETE("/:id", handlers.DeleteDevice)
		devices.POST("/:id/heartbeat", handlers.SendHeartbeat)
		devices.GET("/:id/heartbeat", handlers.GetLatestHeartbeat)
	}

	dashboard := r.Group("/dashboard", middleware.AuthRequired())
	{
		dashboard.GET("/summary", handlers.GetDashboardSummary)
	}

	// WebSocket
	r.GET("/ws", ws.ServeWS)

	log.Println("Server starting on port " + config.App.ServerPort)
	r.Run(":" + config.App.ServerPort)
}