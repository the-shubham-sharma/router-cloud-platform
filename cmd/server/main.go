package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"router-cloud-platform/internal/alert"
	"router-cloud-platform/internal/cache"
	"router-cloud-platform/internal/config"
	"router-cloud-platform/internal/database"
	"router-cloud-platform/internal/handlers"
	"router-cloud-platform/internal/middleware"
	"router-cloud-platform/internal/queue"
	"router-cloud-platform/internal/worker"
	ws "router-cloud-platform/internal/websocket"
)

func main() {
	config.Load()
	database.Connect()
	database.AutoMigrate()
	cache.Connect()
	queue.Connect()
	queue.StartConsumer()

	go ws.GlobalHub.Run()

	worker.Pool = worker.NewHeartbeatWorkerPool(5, 100)
	worker.Pool.Start()

	alert.StartOfflineDetector()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.RateLimit())
	r.Use(middleware.PrometheusMetrics())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

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

	admin := r.Group("/admin", middleware.AuthRequired(), middleware.AdminOnly())
	{
		admin.GET("/devices", handlers.AdminGetAllDevices)
		admin.GET("/users", handlers.AdminGetAllUsers)
		admin.PUT("/users/:id/promote", handlers.AdminPromoteUser)
	}

	r.GET("/ws", ws.ServeWS)

	srv := &http.Server{
		Addr:    ":" + config.App.ServerPort,
		Handler: r,
	}

	go func() {
		log.Println("Server starting on port " + config.App.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	queue.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Forced shutdown: %v", err)
	}

	log.Println("Server stopped cleanly")
}