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

func main() {
	brokerURL := "tcp://broker.hivemq.com:1883"
	clientID := "esp32_simulator_001"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", token.Error())
	}

	fmt.Println("ESP32 Simulator connected to MQTT broker")
	fmt.Println("Simulating sensor data from Sugar Cane plantation...")

	devices := []struct {
		ID       string
		FarmName string
		Lat      float64
		Lng      float64
	}{
		{"ESP32_001", "Farm_A", -7.250445, 112.768845},
		{"ESP32_002", "Farm_B", -7.251234, 112.769876},
		{"ESP32_003", "Farm_C", -7.252567, 112.771234},
	}

	for _, device := range devices {
		sendDeviceStatus(client, device.ID, device.FarmName, true)
		time.Sleep(1 * time.Second)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for _, device := range devices {
			sendSensorData(client, device.ID, device.FarmName, device.Lat, device.Lng)
			time.Sleep(2 * time.Second)
		}
	}
}

func sendSensorData(client mqtt.Client, deviceID, farmName string, lat, lng float64) {
	payload := SensorPayload{
		DeviceID:        deviceID,
		FarmName:        farmName,
		Nitrogen:        generateRealisticValue(15, 35, 25),
		Phosphorus:      generateRealisticValue(10, 25, 18),
		Potassium:       generateRealisticValue(120, 200, 160),
		Temperature:     generateRealisticValue(25, 32, 28),
		Humidity:        generateRealisticValue(60, 80, 70),
		pH:              generateRealisticValue(5.5, 8.5, 6.8),
		Latitude:        lat + (rand.Float64()-0.5)*0.001,
		Longitude:       lng + (rand.Float64()-0.5)*0.001,
		Location:        fmt.Sprintf("Section %s-%d", farmName[len(farmName)-1:], rand.Intn(5)+1),
		BatteryLevel:    float64(rand.Intn(40) + 60),
		SignalStrength:  rand.Intn(30) - 85,
		FirmwareVersion: "v1.2.3",
		Timestamp:       time.Now().Unix(),
	}

	if rand.Float64() < 0.1 {
		switch rand.Intn(4) {
		case 0:
			payload.Nitrogen = 15
		case 1:
			payload.Phosphorus = 12
		case 2:
			payload.Potassium = 130
		case 3:
			payload.pH = 5.2
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
		deviceID, farmName, payload.Nitrogen, payload.Phosphorus, payload.Potassium, payload.pH)
}

func sendDeviceStatus(client mqtt.Client, deviceID, farmName string, isOnline bool) {
	payload := StatusPayload{
		DeviceID:        deviceID,
		FarmName:        farmName,
		IsOnline:        isOnline,
		BatteryLevel:    float64(rand.Intn(40) + 60),
		SignalStrength:  rand.Intn(30) - 85,
		FirmwareVersion: "v1.2.3",
		Location:        fmt.Sprintf("Section %s-%d", farmName[len(farmName)-1:], rand.Intn(5)+1),
		Timestamp:       time.Now().Unix(),
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
	fmt.Printf("📡 Device %s (%s): %s\n", deviceID, farmName, status)
}

func generateRealisticValue(min, max, center float64) float64 {
	variance := (max - min) / 6
	value := rand.NormFloat64()*variance + center

	if value < min {
		value = min
	} else if value > max {
		value = max
	}

	return math.Round(value*100) / 100
}
