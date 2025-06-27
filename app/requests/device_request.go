package requests

// DeviceCommandRequest represents the request body for sending commands to IoT devices
// Example JSON:
// {
//   "command": "SET_THRESHOLD",
//   "payload": {
//     "nitrogen_min": 40.0,
//     "nitrogen_max": 80.0,
//     "temperature_max": 35.0
//   }
// }
type DeviceCommandRequest struct {
	Command string      `json:"command" binding:"required" example:"SET_THRESHOLD" swaggertype:"string"`
	Payload interface{} `json:"payload" example:"{\"nitrogen_min\":40.0,\"nitrogen_max\":80.0,\"temperature_max\":35.0}" swaggertype:"object"`
}
