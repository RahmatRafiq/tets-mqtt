package requests

// SensorDataFilterRequest represents the request body for filtering sensor data
// Example JSON:
// {
//   "device_id": "ESP32_001",
//   "farm_name": "Farm_A",
//   "start_date": "2025-06-01",
//   "end_date": "2025-06-27",
//   "limit": 50,
//   "offset": 0
// }
type SensorDataFilterRequest struct {
	DeviceID  string `json:"device_id" form:"device_id" example:"ESP32_001" swaggertype:"string"`
	FarmName  string `json:"farm_name" form:"farm_name" example:"Farm_A" swaggertype:"string"`
	StartDate string `json:"start_date" form:"start_date" example:"2025-06-01" swaggertype:"string"`
	EndDate   string `json:"end_date" form:"end_date" example:"2025-06-27" swaggertype:"string"`
	Limit     int    `json:"limit" form:"limit" example:"50" default:"50"`
	Offset    int    `json:"offset" form:"offset" example:"0" default:"0"`
}
