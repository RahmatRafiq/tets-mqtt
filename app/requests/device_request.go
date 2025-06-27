package requests

type DeviceCommandRequest struct {
	Command string      `json:"command" binding:"required"`
	Payload interface{} `json:"payload"`
}
