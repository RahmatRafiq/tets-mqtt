package models

import (
	"time"
)

type SensorAlert struct {
	ID             uint       `json:"id" gorm:"primarykey"`
	DeviceID       string     `json:"device_id" gorm:"not null;index"`
	FarmName       string     `json:"farm_name" gorm:"not null"`
	AlertType      string     `json:"alert_type" gorm:"not null"`
	Message        string     `json:"message" gorm:"not null"`
	Severity       string     `json:"severity" gorm:"not null"`
	SensorValue    *float64   `json:"sensor_value"`
	ThresholdValue *float64   `json:"threshold_value"`
	IsResolved     bool       `json:"is_resolved" gorm:"default:false"`
	ResolvedAt     *time.Time `json:"resolved_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (SensorAlert) TableName() string {
	return "sensor_alerts"
}
