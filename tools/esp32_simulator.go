package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SensorPayload represents the structure of data sent from ESP32 sensors
type SensorPayload struct {
	DeviceID       string  `json:"device_id"`
	KebunName      string  `json:"kebun_name"`
	Nitrogen       float64 `json:"nitrogen"`
	Phosphorus     float64 `json:"phosphorus"`
	Potassium      float64 `json:"potassium"`
	Temperature    float64 `json:"temperature"`
	Humidity       float64 `json:"humidity"`
	pH             float64 `json:"ph"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	BatteryLevel   int     `json:"battery_level"`
	SignalStrength int     `json:"signal_strength"`
	Timestamp      int64   `json:"timestamp"`
}

// StatusPayload represents device status updates
type StatusPayload struct {
	DeviceID       string `json:"device_id"`
	KebunName      string `json:"kebun_name"`
	IsOnline       bool   `json:"is_online"`
	BatteryLevel   int    `json:"battery_level"`
	SignalStrength int    `json:"signal_strength"`
	Timestamp      int64  `json:"timestamp"`
}

func main() {
	// MQTT broker configuration
	brokerURL := "tcp://broker.hivemq.com:1883"
	clientID := "esp32_simulator_001"

	// Create MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	// Create and connect client
	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	fmt.Println("ESP32 Simulator connected to MQTT broker")
	fmt.Println("Simulating sensor data from Sugar Cane plantation...")

	// Simulate multiple devices
	devices := []struct {
		ID        string
		KebunName string
		Lat       float64
		Lng       float64
	}{
		{"ESP32_001", "Kebun A", -7.250445, 112.768845},
		{"ESP32_002", "Kebun B", -7.251234, 112.769876},
		{"ESP32_003", "Kebun C", -7.252567, 112.771234},
	}

	// Send initial device status
	for _, device := range devices {
		sendDeviceStatus(client, device.ID, device.KebunName, true)
		time.Sleep(1 * time.Second)
	}

	// Continuously send sensor data
	ticker := time.NewTicker(30 * time.Second) // Send data every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, device := range devices {
				sendSensorData(client, device.ID, device.KebunName, device.Lat, device.Lng)
				time.Sleep(2 * time.Second) // Stagger the sends
			}
		}
	}
}

func sendSensorData(client mqtt.Client, deviceID, kebunName string, lat, lng float64) {
	// Generate realistic sensor data for sugar cane
	payload := SensorPayload{
		DeviceID:       deviceID,
		KebunName:      kebunName,
		Nitrogen:       generateRealisticValue(15, 35, 25),    // Optimal: 20-30 mg/kg
		Phosphorus:     generateRealisticValue(10, 25, 18),    // Optimal: 15-20 mg/kg
		Potassium:      generateRealisticValue(120, 200, 160), // Optimal: 150-180 mg/kg
		Temperature:    generateRealisticValue(25, 32, 28),    // Tropical temperature
		Humidity:       generateRealisticValue(60, 80, 70),    // Soil humidity %
		pH:             generateRealisticValue(5.5, 8.5, 6.8), // Optimal: 6.0-8.0
		Latitude:       lat + (rand.Float64()-0.5)*0.001,      // Small variations
		Longitude:      lng + (rand.Float64()-0.5)*0.001,      // Small variations
		BatteryLevel:   rand.Intn(40) + 60,                    // 60-100%
		SignalStrength: rand.Intn(30) - 85,                    // -85 to -55 dBm
		Timestamp:      time.Now().Unix(),
	}

	// Occasionally generate alert conditions
	if rand.Float64() < 0.1 { // 10% chance
		switch rand.Intn(4) {
		case 0:
			payload.Nitrogen = 15 // Low nitrogen
		case 1:
			payload.Phosphorus = 12 // Low phosphorus
		case 2:
			payload.Potassium = 130 // Low potassium
		case 3:
			payload.pH = 5.2 // Low pH
		}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling sensor data: %v", err)
		return
	}

	topic := fmt.Sprintf("sugar_vestrack/sensor/%s/data", deviceID)
	token := client.Publish(topic, 1, false, jsonData)
	if token.Wait() && token.Error() != nil {
		log.Printf("Error publishing sensor data: %v", token.Error())
		return
	}

	fmt.Printf("📊 Sent sensor data from %s (%s): N=%.1f, P=%.1f, K=%.1f, pH=%.1f\n",
		deviceID, kebunName, payload.Nitrogen, payload.Phosphorus, payload.Potassium, payload.pH)
}

func sendDeviceStatus(client mqtt.Client, deviceID, kebunName string, isOnline bool) {
	payload := StatusPayload{
		DeviceID:       deviceID,
		KebunName:      kebunName,
		IsOnline:       isOnline,
		BatteryLevel:   rand.Intn(40) + 60, // 60-100%
		SignalStrength: rand.Intn(30) - 85, // -85 to -55 dBm
		Timestamp:      time.Now().Unix(),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling status data: %v", err)
		return
	}

	topic := fmt.Sprintf("sugar_vestrack/device/%s/status", deviceID)
	token := client.Publish(topic, 1, false, jsonData)
	if token.Wait() && token.Error() != nil {
		log.Printf("Error publishing status data: %v", token.Error())
		return
	}

	status := "🟢 Online"
	if !isOnline {
		status = "🔴 Offline"
	}
	fmt.Printf("📡 Device %s (%s): %s\n", deviceID, kebunName, status)
}

func generateRealisticValue(min, max, center float64) float64 {
	// Generate values with normal distribution around center
	variance := (max - min) / 6 // 99.7% of values within range
	value := rand.NormFloat64()*variance + center

	// Clamp to min/max range
	if value < min {
		value = min
	} else if value > max {
		value = max
	}

	return math.Round(value*100) / 100 // Round to 2 decimal places
}
