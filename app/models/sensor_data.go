package models

import (
	"time"

	"gorm.io/gorm"
)

// SensorData represents NPK sensor readings from ESP32 devices
type SensorData struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	DeviceID    string         `json:"device_id" gorm:"not null;index"`
	FarmName    string         `json:"farm_name" gorm:"not null"`
	Nitrogen    float64        `json:"nitrogen"`
	Phosphorus  float64        `json:"phosphorus"`
	Potassium   float64        `json:"potassium"`
	Temperature float64        `json:"temperature"`
	Humidity    float64        `json:"humidity"`
	PH          float64        `json:"ph"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	Location    string         `json:"location"`
	Timestamp   time.Time      `json:"timestamp" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (SensorData) TableName() string {
	return "sensor_data"
}
