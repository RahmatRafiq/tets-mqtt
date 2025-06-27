package services

type SensorPayload struct {
	DeviceID        string  `json:"device_id"`
	FarmName        string  `json:"farm_name"`
	Nitrogen        float64 `json:"nitrogen"`
	Phosphorus      float64 `json:"phosphorus"`
	Potassium       float64 `json:"potassium"`
	Temperature     float64 `json:"temperature"`
	Humidity        float64 `json:"humidity"`
	pH              float64 `json:"ph"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Location        string  `json:"location"`
	BatteryLevel    float64 `json:"battery_level"`
	SignalStrength  int     `json:"signal_strength"`
	FirmwareVersion string  `json:"firmware_version"`
	Timestamp       int64   `json:"timestamp"`
}

type StatusPayload struct {
	DeviceID        string  `json:"device_id"`
	FarmName        string  `json:"farm_name"`
	IsOnline        bool    `json:"is_online"`
	BatteryLevel    float64 `json:"battery_level"`
	SignalStrength  int     `json:"signal_strength"`
	FirmwareVersion string  `json:"firmware_version"`
	Location        string  `json:"location"`
	Timestamp       int64   `json:"timestamp"`
}
