package requests

// DeviceCommandRequest represents a command to be sent to an IoT device
type DeviceCommandRequest struct {
	Command string      `json:"command" binding:"required"` // Command type: "restart", "calibrate", "set_interval", etc.
	Payload interface{} `json:"payload"`                    // Command-specific payload
}

// SensorDataFilterRequest represents filters for sensor data queries
type SensorDataFilterRequest struct {
	DeviceID  string `json:"device_id" form:"device_id"`
	KebunName string `json:"kebun_name" form:"kebun_name"`
	StartDate string `json:"start_date" form:"start_date"` // YYYY-MM-DD format
	EndDate   string `json:"end_date" form:"end_date"`     // YYYY-MM-DD format
	Limit     int    `json:"limit" form:"limit"`
	Offset    int    `json:"offset" form:"offset"`
}
