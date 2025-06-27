package cmd

import (
	"log"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"
)

func FixSensorMigration() {
	// Connect to database
	if facades.DB == nil {
		log.Fatal("Database connection not available")
	}

	log.Println("Starting sensor tables migration fix...")

	// Drop existing tables to avoid migration conflicts
	log.Println("Dropping existing sensor tables...")
	if err := facades.DB.Migrator().DropTable(&models.SensorData{}); err != nil {
		log.Printf("Warning: Failed to drop sensor_data table: %v", err)
	}

	if err := facades.DB.Migrator().DropTable(&models.DeviceStatus{}); err != nil {
		log.Printf("Warning: Failed to drop device_status table: %v", err)
	}

	if err := facades.DB.Migrator().DropTable(&models.SensorAlert{}); err != nil {
		log.Printf("Warning: Failed to drop sensor_alerts table: %v", err)
	}

	log.Println("Creating new sensor tables with correct constraints...")
	
	// Create tables with new structure
	if err := facades.DB.AutoMigrate(
		&models.SensorData{},
		&models.DeviceStatus{},
		&models.SensorAlert{},
	); err != nil {
		log.Fatalf("Failed to migrate sensor tables: %v", err)
	}

	log.Println("Sensor tables migration completed successfully!")
}
