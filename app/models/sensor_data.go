package models

import (
	"time"

	"gorm.io/gorm"
)

// SensorData represents NPK sensor readings from ESP32 devices
type SensorData struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	DeviceID    string         `json:"device_id" gorm:"not null;index"`
	KebunName   string         `json:"kebun_name" gorm:"not null"`
	Nitrogen    float64        `json:"nitrogen"`    // N value
	Phosphorus  float64        `json:"phosphorus"`  // P value
	Potassium   float64        `json:"potassium"`   // K value
	Temperature float64        `json:"temperature"` // Soil temperature
	Humidity    float64        `json:"humidity"`    // Soil humidity
	PH          float64        `json:"ph"`          // Soil pH level
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	Timestamp   time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName sets the table name for GORM
func (SensorData) TableName() string {
	return "sensor_data"
}

// DeviceStatus represents the online/offline status of IoT devices
type DeviceStatus struct {
	ID             uint      `json:"id" gorm:"primarykey"`
	DeviceID       string    `json:"device_id" gorm:"uniqueIndex;not null"`
	KebunName      string    `json:"kebun_name" gorm:"not null"`
	IsOnline       bool      `json:"is_online" gorm:"default:false"`
	LastSeen       time.Time `json:"last_seen"`
	BatteryLevel   int       `json:"battery_level"`   // Battery percentage
	SignalStrength int       `json:"signal_strength"` // LoRa signal strength
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName sets the table name for GORM
func (DeviceStatus) TableName() string {
	return "device_status"
}

// SensorAlert represents alerts based on sensor thresholds
type SensorAlert struct {
	ID         uint       `json:"id" gorm:"primarykey"`
	DeviceID   string     `json:"device_id" gorm:"not null;index"`
	KebunName  string     `json:"kebun_name" gorm:"not null"`
	AlertType  string     `json:"alert_type" gorm:"not null"` // "npk_low", "ph_abnormal", "moisture_low", etc.
	Message    string     `json:"message" gorm:"not null"`
	Severity   string     `json:"severity" gorm:"not null"` // "low", "medium", "high", "critical"
	IsResolved bool       `json:"is_resolved" gorm:"default:false"`
	ResolvedAt *time.Time `json:"resolved_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// TableName sets the table name for GORM
func (SensorAlert) TableName() string {
	return "sensor_alerts"
}
