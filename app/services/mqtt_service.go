package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTService struct {
	client      mqtt.Client
	isConnected bool
	brokerURL   string
	clientID    string
	username    string
	password    string
}

// SensorPayload represents the structure of data received from ESP32 sensors
type SensorPayload struct {
	DeviceID         string  `json:"device_id"`
	FarmName         string  `json:"farm_name"`
	Nitrogen         float64 `json:"nitrogen"`
	Phosphorus       float64 `json:"phosphorus"`
	Potassium        float64 `json:"potassium"`
	Temperature      float64 `json:"temperature"`
	Humidity         float64 `json:"humidity"`
	pH               float64 `json:"ph"`
	Latitude         float64 `json:"latitude"`
	Longitude        float64 `json:"longitude"`
	Location         string  `json:"location"`
	BatteryLevel     float64 `json:"battery_level"`
	SignalStrength   int     `json:"signal_strength"`
	FirmwareVersion  string  `json:"firmware_version"`
	Timestamp        int64   `json:"timestamp"`
}

// StatusPayload represents device status updates
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

var mqttServiceInstance *MQTTService

// NewMQTTService creates a new MQTT service instance
func NewMQTTService(brokerURL, clientID, username, password string) *MQTTService {
	if mqttServiceInstance == nil {
		mqttServiceInstance = &MQTTService{
			brokerURL: brokerURL,
			clientID:  clientID,
			username:  username,
			password:  password,
		}
	}
	return mqttServiceInstance
}

// GetMQTTService returns the singleton instance
func GetMQTTService() *MQTTService {
	return mqttServiceInstance
}

// Connect establishes connection to MQTT broker
func (m *MQTTService) Connect() error {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(m.brokerURL)
	opts.SetClientID(m.clientID)
	opts.SetUsername(m.username)
	opts.SetPassword(m.password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)

	// Set connection handlers
	opts.SetOnConnectHandler(m.onConnect)
	opts.SetConnectionLostHandler(m.onConnectionLost)

	m.client = mqtt.NewClient(opts)

	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	m.isConnected = true
	log.Println("Successfully connected to MQTT broker")
	return nil
}

// onConnect callback when connected to broker
func (m *MQTTService) onConnect(client mqtt.Client) {
	log.Println("MQTT client connected to broker")

	// Subscribe to sensor data topics
	m.subscribeToTopics()
}

// onConnectionLost callback when connection is lost
func (m *MQTTService) onConnectionLost(client mqtt.Client, err error) {
	log.Printf("MQTT connection lost: %v", err)
	m.isConnected = false
}

// subscribeToTopics subscribes to all necessary MQTT topics
func (m *MQTTService) subscribeToTopics() {
	topics := map[string]mqtt.MessageHandler{
		"sugar_vestrack/sensor/+/data":   m.handleSensorData,
		"sugar_vestrack/device/+/status": m.handleDeviceStatus,
		"sugar_vestrack/device/+/alert":  m.handleDeviceAlert,
	}

	for topic, handler := range topics {
		token := m.client.Subscribe(topic, 1, handler)
		if token.Wait() && token.Error() != nil {
			log.Printf("Failed to subscribe to topic %s: %v", topic, token.Error())
		} else {
			log.Printf("Successfully subscribed to topic: %s", topic)
		}
	}
}

// handleSensorData processes incoming sensor data
func (m *MQTTService) handleSensorData(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received sensor data from topic: %s", msg.Topic())

	var payload SensorPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("Error parsing sensor data: %v", err)
		return
	}

	// Save sensor data to database
	sensorData := models.SensorData{
		DeviceID:    payload.DeviceID,
		FarmName:    payload.FarmName,
		Nitrogen:    payload.Nitrogen,
		Phosphorus:  payload.Phosphorus,
		Potassium:   payload.Potassium,
		Temperature: payload.Temperature,
		Humidity:    payload.Humidity,
		PH:          payload.pH,
		Latitude:    payload.Latitude,
		Longitude:   payload.Longitude,
		Location:    payload.Location,
		Timestamp:   time.Unix(payload.Timestamp, 0),
	}

	db := facades.DB
	if err := db.Create(&sensorData).Error; err != nil {
		log.Printf("Error saving sensor data: %v", err)
		return
	}

	// Update device status
	m.updateDeviceStatus(payload.DeviceID, payload.FarmName, true, payload.BatteryLevel, payload.SignalStrength, payload.FirmwareVersion, payload.Location)

	// Check for alerts
	m.checkSensorAlerts(&sensorData)

	log.Printf("Sensor data saved for device %s in %s", payload.DeviceID, payload.FarmName)
}

// handleDeviceStatus processes device status updates
func (m *MQTTService) handleDeviceStatus(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received device status from topic: %s", msg.Topic())

	var payload StatusPayload
	if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
		log.Printf("Error parsing device status: %v", err)
		return
	}

	m.updateDeviceStatus(payload.DeviceID, payload.FarmName, payload.IsOnline, payload.BatteryLevel, payload.SignalStrength, payload.FirmwareVersion, payload.Location)
}

// handleDeviceAlert processes device alerts
func (m *MQTTService) handleDeviceAlert(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received device alert from topic: %s", msg.Topic())

	// Handle custom alerts from devices
	// Implementation depends on alert format from ESP32
}

// updateDeviceStatus updates or creates device status record
func (m *MQTTService) updateDeviceStatus(deviceID, farmName string, isOnline bool, batteryLevel float64, signalStrength int, firmwareVersion, location string) {
	db := facades.DB

	var deviceStatus models.DeviceStatus
	result := db.Where("device_id = ?", deviceID).First(&deviceStatus)

	if result.Error != nil {
		// Create new device status
		deviceStatus = models.DeviceStatus{
			DeviceID:        deviceID,
			FarmName:        farmName,
			IsOnline:        isOnline,
			LastSeen:        time.Now(),
			BatteryLevel:    batteryLevel,
			SignalStrength:  signalStrength,
			FirmwareVersion: firmwareVersion,
			Location:        location,
		}
		db.Create(&deviceStatus)
	} else {
		// Update existing device status
		db.Model(&deviceStatus).Updates(models.DeviceStatus{
			IsOnline:        isOnline,
			LastSeen:        time.Now(),
			BatteryLevel:    batteryLevel,
			SignalStrength:  signalStrength,
			FirmwareVersion: firmwareVersion,
			Location:        location,
		})
	}
}

// checkSensorAlerts checks sensor values and creates alerts if needed
func (m *MQTTService) checkSensorAlerts(data *models.SensorData) {
	db := facades.DB

	alerts := []models.SensorAlert{}

	// Check NPK levels (example thresholds)
	if data.Nitrogen < 20 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "nitrogen_low",
			Message:        fmt.Sprintf("Nitrogen level too low: %.2f mg/kg", data.Nitrogen),
			Severity:       "medium",
			SensorValue:    &data.Nitrogen,
			ThresholdValue: func() *float64 { v := 20.0; return &v }(),
		})
	}

	if data.Phosphorus < 15 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "phosphorus_low",
			Message:        fmt.Sprintf("Phosphorus level too low: %.2f mg/kg", data.Phosphorus),
			Severity:       "medium",
			SensorValue:    &data.Phosphorus,
			ThresholdValue: func() *float64 { v := 15.0; return &v }(),
		})
	}

	if data.Potassium < 150 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "potassium_low",
			Message:        fmt.Sprintf("Potassium level too low: %.2f mg/kg", data.Potassium),
			Severity:       "medium",
			SensorValue:    &data.Potassium,
			ThresholdValue: func() *float64 { v := 150.0; return &v }(),
		})
	}

	// Check pH levels (optimal range for sugar cane: 6.0-8.0)
	if data.PH < 6.0 || data.PH > 8.0 {
		severity := "low"
		if data.PH < 5.5 || data.PH > 8.5 {
			severity = "high"
		}
		thresholdValue := 6.0
		if data.PH > 8.0 {
			thresholdValue = 8.0
		}
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "ph_abnormal",
			Message:        fmt.Sprintf("pH level abnormal: %.2f", data.PH),
			Severity:       severity,
			SensorValue:    &data.PH,
			ThresholdValue: &thresholdValue,
		})
	}

	// Save alerts to database
	for _, alert := range alerts {
		db.Create(&alert)
	}
}

// PublishCommand publishes a command to a specific device
func (m *MQTTService) PublishCommand(deviceID, command string, payload interface{}) error {
	if !m.isConnected {
		return fmt.Errorf("MQTT client not connected")
	}

	topic := fmt.Sprintf("sugar_vestrack/device/%s/command", deviceID)

	data := map[string]interface{}{
		"command":   command,
		"payload":   payload,
		"timestamp": time.Now().Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling command: %v", err)
	}

	token := m.client.Publish(topic, 1, false, jsonData)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("error publishing command: %v", token.Error())
	}

	return nil
}

// IsConnected returns the connection status
func (m *MQTTService) IsConnected() bool {
	return m.isConnected
}

// Disconnect closes the MQTT connection
func (m *MQTTService) Disconnect() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
		m.isConnected = false
		log.Println("MQTT client disconnected")
	}
}
