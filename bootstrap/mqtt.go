package bootstrap

import (
	"log"
	"os"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"
)

// InitializeMQTT initializes the MQTT service and auto-migrates tables
func InitializeMQTT() {
	// Auto-migrate sensor tables
	db := facades.DB
	if err := db.AutoMigrate(
		&models.SensorData{},
		&models.DeviceStatus{},
		&models.SensorAlert{},
	); err != nil {
		log.Printf("Failed to migrate sensor tables: %v", err)
	} else {
		log.Println("Sensor tables migrated successfully")
	}

	// Initialize MQTT service
	brokerURL := helpers.GetEnv("MQTT_BROKER_URL", "tcp://localhost:1883")
	clientID := helpers.GetEnv("MQTT_CLIENT_ID", "sugar_vestrack_server")
	username := helpers.GetEnv("MQTT_USERNAME", "")
	password := helpers.GetEnv("MQTT_PASSWORD", "")

	// Create MQTT service
	mqttService := services.NewMQTTService(brokerURL, clientID, username, password)

	// Connect to MQTT broker
	if err := mqttService.Connect(); err != nil {
		log.Printf("Failed to connect to MQTT broker: %v", err)
		log.Println("MQTT functionality will be disabled")
	} else {
		log.Println("MQTT service initialized successfully")
	}
}

// InitializeEnvironment loads environment variables for MQTT
func InitializeEnvironment() {
	// Set default MQTT environment variables if not present
	if os.Getenv("MQTT_BROKER_URL") == "" {
		os.Setenv("MQTT_BROKER_URL", "tcp://broker.hivemq.com:1883")
	}

	if os.Getenv("MQTT_CLIENT_ID") == "" {
		os.Setenv("MQTT_CLIENT_ID", "sugar_vestrack_server")
	}
}
