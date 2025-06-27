package models

import (
	"time"
)

type DeviceStatus struct {
	ID              uint      `json:"id" gorm:"primarykey"`
	DeviceID        string    `json:"device_id" gorm:"type:varchar(100);uniqueIndex;not null"`
	FarmName        string    `json:"farm_name" gorm:"type:varchar(255);not null"`
	IsOnline        bool      `json:"is_online" gorm:"default:false"`
	LastSeen        time.Time `json:"last_seen"`
	BatteryLevel    float64   `json:"battery_level"`
	SignalStrength  int       `json:"signal_strength"`
	FirmwareVersion string    `json:"firmware_version" gorm:"type:varchar(50)"`
	Location        string    `json:"location" gorm:"type:varchar(500)"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (DeviceStatus) TableName() string {
	return "device_status"
}
