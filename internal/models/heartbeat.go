package models

import (
	"time"
	"github.com/google/uuid"
)

type Heartbeat struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DeviceID  uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	CPU       float64   `json:"cpu"`
	Memory    float64   `json:"memory"`
	Latency   float64   `json:"latency"`
	Bandwidth float64   `json:"bandwidth"`
	CreatedAt time.Time `json:"created_at"`
	Device    Device    `gorm:"foreignKey:DeviceID" json:"-"`
}