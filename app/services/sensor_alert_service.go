package services

import (
	"fmt"

	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/facades"
)

type SensorAlertService struct{}

func NewSensorAlertService() *SensorAlertService {
	return &SensorAlertService{}
}

func (s *SensorAlertService) CheckSensorAlerts(data *models.SensorData) {
	db := facades.DB

	alerts := []models.SensorAlert{}

	alerts = append(alerts, s.checkNitrogenLevel(data)...)
	alerts = append(alerts, s.checkPhosphorusLevel(data)...)
	alerts = append(alerts, s.checkPotassiumLevel(data)...)
	alerts = append(alerts, s.checkPHLevel(data)...)
	alerts = append(alerts, s.checkTemperatureLevel(data)...)
	alerts = append(alerts, s.checkHumidityLevel(data)...)

	for _, alert := range alerts {
		db.Create(&alert)
	}
}

func (s *SensorAlertService) checkNitrogenLevel(data *models.SensorData) []models.SensorAlert {
	var alerts []models.SensorAlert

	if data.Nitrogen < 20 {
		severity := "medium"
		if data.Nitrogen < 15 {
			severity = "high"
		}
		if data.Nitrogen < 10 {
			severity = "critical"
		}

		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "nitrogen_low",
			Message:        fmt.Sprintf("Nitrogen level too low: %.2f mg/kg (optimal: 20-30 mg/kg)", data.Nitrogen),
			Severity:       severity,
			SensorValue:    &data.Nitrogen,
			ThresholdValue: func() *float64 { v := 20.0; return &v }(),
		})
	} else if data.Nitrogen > 40 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "nitrogen_high",
			Message:        fmt.Sprintf("Nitrogen level too high: %.2f mg/kg (optimal: 20-30 mg/kg)", data.Nitrogen),
			Severity:       "medium",
			SensorValue:    &data.Nitrogen,
			ThresholdValue: func() *float64 { v := 40.0; return &v }(),
		})
	}

	return alerts
}

func (s *SensorAlertService) checkPhosphorusLevel(data *models.SensorData) []models.SensorAlert {
	var alerts []models.SensorAlert

	if data.Phosphorus < 15 {
		severity := "medium"
		if data.Phosphorus < 10 {
			severity = "high"
		}
		if data.Phosphorus < 5 {
			severity = "critical"
		}

		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "phosphorus_low",
			Message:        fmt.Sprintf("Phosphorus level too low: %.2f mg/kg (optimal: 15-20 mg/kg)", data.Phosphorus),
			Severity:       severity,
			SensorValue:    &data.Phosphorus,
			ThresholdValue: func() *float64 { v := 15.0; return &v }(),
		})
	} else if data.Phosphorus > 30 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "phosphorus_high",
			Message:        fmt.Sprintf("Phosphorus level too high: %.2f mg/kg (optimal: 15-20 mg/kg)", data.Phosphorus),
			Severity:       "medium",
			SensorValue:    &data.Phosphorus,
			ThresholdValue: func() *float64 { v := 30.0; return &v }(),
		})
	}

	return alerts
}

func (s *SensorAlertService) checkPotassiumLevel(data *models.SensorData) []models.SensorAlert {
	var alerts []models.SensorAlert

	if data.Potassium < 150 {
		severity := "medium"
		if data.Potassium < 120 {
			severity = "high"
		}
		if data.Potassium < 100 {
			severity = "critical"
		}

		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "potassium_low",
			Message:        fmt.Sprintf("Potassium level too low: %.2f mg/kg (optimal: 150-180 mg/kg)", data.Potassium),
			Severity:       severity,
			SensorValue:    &data.Potassium,
			ThresholdValue: func() *float64 { v := 150.0; return &v }(),
		})
	} else if data.Potassium > 220 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "potassium_high",
			Message:        fmt.Sprintf("Potassium level too high: %.2f mg/kg (optimal: 150-180 mg/kg)", data.Potassium),
			Severity:       "medium",
			SensorValue:    &data.Potassium,
			ThresholdValue: func() *float64 { v := 220.0; return &v }(),
		})
	}

	return alerts
}

func (s *SensorAlertService) checkPHLevel(data *models.SensorData) []models.SensorAlert {
	var alerts []models.SensorAlert

	if data.PH < 6.0 || data.PH > 8.0 {
		severity := "low"
		thresholdValue := 6.0
		message := fmt.Sprintf("pH level too low: %.2f (optimal: 6.0-8.0)", data.PH)

		if data.PH > 8.0 {
			thresholdValue = 8.0
			message = fmt.Sprintf("pH level too high: %.2f (optimal: 6.0-8.0)", data.PH)
		}

		if data.PH < 5.5 || data.PH > 8.5 {
			severity = "high"
		}
		if data.PH < 5.0 || data.PH > 9.0 {
			severity = "critical"
		}

		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "ph_abnormal",
			Message:        message,
			Severity:       severity,
			SensorValue:    &data.PH,
			ThresholdValue: &thresholdValue,
		})
	}

	return alerts
}

func (s *SensorAlertService) checkTemperatureLevel(data *models.SensorData) []models.SensorAlert {
	var alerts []models.SensorAlert

	if data.Temperature < 20 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "temperature_low",
			Message:        fmt.Sprintf("Temperature too low: %.2f°C (optimal: 25-35°C)", data.Temperature),
			Severity:       "medium",
			SensorValue:    &data.Temperature,
			ThresholdValue: func() *float64 { v := 20.0; return &v }(),
		})
	} else if data.Temperature > 40 {
		severity := "medium"
		if data.Temperature > 45 {
			severity = "high"
		}

		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "temperature_high",
			Message:        fmt.Sprintf("Temperature too high: %.2f°C (optimal: 25-35°C)", data.Temperature),
			Severity:       severity,
			SensorValue:    &data.Temperature,
			ThresholdValue: func() *float64 { v := 40.0; return &v }(),
		})
	}

	return alerts
}

func (s *SensorAlertService) checkHumidityLevel(data *models.SensorData) []models.SensorAlert {
	var alerts []models.SensorAlert

	if data.Humidity < 50 {
		severity := "medium"
		if data.Humidity < 40 {
			severity = "high"
		}

		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "humidity_low",
			Message:        fmt.Sprintf("Soil humidity too low: %.2f%% (optimal: 60-80%%)", data.Humidity),
			Severity:       severity,
			SensorValue:    &data.Humidity,
			ThresholdValue: func() *float64 { v := 50.0; return &v }(),
		})
	} else if data.Humidity > 90 {
		alerts = append(alerts, models.SensorAlert{
			DeviceID:       data.DeviceID,
			FarmName:       data.FarmName,
			AlertType:      "humidity_high",
			Message:        fmt.Sprintf("Soil humidity too high: %.2f%% (optimal: 60-80%%)", data.Humidity),
			Severity:       "medium",
			SensorValue:    &data.Humidity,
			ThresholdValue: func() *float64 { v := 90.0; return &v }(),
		})
	}

	return alerts
}
