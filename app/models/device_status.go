package models

import (
	"time"
)

type DeviceStatus struct {
	ID              uint      `json:"id" gorm:"primarykey"`
	DeviceID        string    `json:"device_id" gorm:"uniqueIndex;not null"`
	FarmName        string    `json:"farm_name" gorm:"not null"`
	IsOnline        bool      `json:"is_online" gorm:"default:false"`
	LastSeen        time.Time `json:"last_seen"`
	BatteryLevel    float64   `json:"battery_level"`
	SignalStrength  int       `json:"signal_strength"`
	FirmwareVersion string    `json:"firmware_version"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (DeviceStatus) TableName() string {
	return "device_status"
}
