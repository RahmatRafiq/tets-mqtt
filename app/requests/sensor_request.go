package requests

type SensorDataFilterRequest struct {
	DeviceID  string `json:"device_id" form:"device_id"`
	FarmName  string `json:"farm_name" form:"farm_name"`
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
	Limit     int    `json:"limit" form:"limit"`
	Offset    int    `json:"offset" form:"offset"`
}
