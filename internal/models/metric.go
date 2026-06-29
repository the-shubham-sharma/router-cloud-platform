package models

import (
	"time"
	"github.com/google/uuid"
)

type Metric struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DeviceID  uuid.UUID `gorm:"type:uuid;not null" json:"device_id"`
	Key       string    `gorm:"not null" json:"key"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	Device    Device    `gorm:"foreignKey:DeviceID" json:"-"`
}