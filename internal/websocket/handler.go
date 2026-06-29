package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"router-cloud-platform/internal/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

type StatusMessage struct {
	Type      string    `json:"type"`
	DeviceID  string    `json:"device_id,omitempty"`
	Status    string    `json:"status,omitempty"`
	CPU       float64   `json:"cpu,omitempty"`
	Memory    float64   `json:"memory,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func ServeWS(c *gin.Context) {
	// Validate token from query param
	token := c.Query("token")
	if token == "" {
		utils.Error(c, http.StatusUnauthorized, "Token required")
		return
	}

	// Strip Bearer prefix if present
	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := utils.ValidateToken(token)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid token")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: claims.UserID.String(),
	}

	GlobalHub.register <- client

	// Send welcome message
	welcome := StatusMessage{
		Type:      "connected",
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(welcome)
	client.send <- data

	go client.writePump()

	// Keep reading to detect disconnection
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			GlobalHub.unregister <- client
			break
		}
	}
}

func BroadcastDeviceUpdate(deviceID, status string, cpu, memory float64) {
	msg := StatusMessage{
		Type:      "device_update",
		DeviceID:  deviceID,
		Status:    status,
		CPU:       cpu,
		Memory:    memory,
		Timestamp: time.Now(),
	}
	data, _ := json.Marshal(msg)
	GlobalHub.Broadcast(data)
}