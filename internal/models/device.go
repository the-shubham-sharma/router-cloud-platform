package models

import (
	"time"
	"github.com/google/uuid"
)

type DeviceStatus string

const (
	StatusOnline  DeviceStatus = "online"
	StatusOffline DeviceStatus = "offline"
	StatusUnknown DeviceStatus = "unknown"
)

type Device struct {
	ID        uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	Name      string       `gorm:"not null" json:"name"`
	IPAddress string       `gorm:"not null" json:"ip_address"`
	Location  string       `json:"location"`
	Status    DeviceStatus `gorm:"default:unknown" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	User      User         `gorm:"foreignKey:UserID" json:"-"`
}