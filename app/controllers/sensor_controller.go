package controllers

import (
	"net/http"
	"strconv"
	"time"

	"golang_starter_kit_2025/app/helpers"
	"golang_starter_kit_2025/app/models"
	"golang_starter_kit_2025/app/requests"
	"golang_starter_kit_2025/app/services"
	"golang_starter_kit_2025/facades"

	"github.com/gin-gonic/gin"
)

type SensorController struct{}

// GetSensorData gets sensor data with filtering options
// @Summary Get sensor data
// @Description Retrieve sensor data with optional filtering by device, kebun, and time range
// @Tags sensor
// @Accept json
// @Produce json
// @Param device_id query string false "Device ID filter"
// @Param kebun_name query string false "Kebun name filter"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Limit results" default(50)
// @Param offset query int false "Offset results" default(0)
// @Success 200 {object} map[string]interface{}
// @Router /api/sensor/data [get]
func (sc *SensorController) GetSensorData(c *gin.Context) {
	var sensorData []models.SensorData

	query := facades.DB.Model(&models.SensorData{})

	// Apply filters
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if kebunName := c.Query("kebun_name"); kebunName != "" {
		query = query.Where("kebun_name = ?", kebunName)
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("timestamp >= ?", parsedDate)
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("timestamp <= ?", parsedDate.Add(24*time.Hour))
		}
	}

	// Pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	query = query.Limit(limit).Offset(offset).Order("timestamp DESC")

	if err := query.Find(&sensorData).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to retrieve sensor data",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sensor data retrieved successfully",
		"data":    sensorData,
		"count":   len(sensorData),
		"limit":   limit,
		"offset":  offset,
	})
}

// GetLatestSensorData gets the latest sensor readings for each device
// @Summary Get latest sensor data
// @Description Retrieve the most recent sensor readings for all devices or specific device
// @Tags sensor
// @Accept json
// @Produce json
// @Param device_id query string false "Device ID filter"
// @Success 200 {object} map[string]interface{}
// @Router /api/sensor/latest [get]
func (sc *SensorController) GetLatestSensorData(c *gin.Context) {
	deviceID := c.Query("device_id")

	if deviceID != "" {
		// Get latest data for specific device
		var sensorData models.SensorData
		if err := facades.DB.Where("device_id = ?", deviceID).
			Order("timestamp DESC").
			First(&sensorData).Error; err != nil {
			helpers.ResponseError(c, &helpers.ResponseParams[any]{
				Message: "No data found for device",
				Errors:  map[string]string{"error": err.Error()},
			}, http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Latest sensor data retrieved",
			"data":    sensorData,
		})
		return
	}

	// Get latest data for each device
	var sensorData []models.SensorData
	subQuery := facades.DB.Model(&models.SensorData{}).
		Select("device_id, MAX(timestamp) as max_timestamp").
		Group("device_id")

	if err := facades.DB.Table("sensor_data").
		Joins("INNER JOIN (?) as latest ON sensor_data.device_id = latest.device_id AND sensor_data.timestamp = latest.max_timestamp", subQuery).
		Find(&sensorData).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to retrieve latest sensor data",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Latest sensor data retrieved successfully",
		"data":    sensorData,
		"count":   len(sensorData),
	})
}

// GetDeviceStatus gets device status information
func (sc *SensorController) GetDeviceStatus(c *gin.Context) {
	var deviceStatus []models.DeviceStatus

	query := facades.DB.Model(&models.DeviceStatus{})

	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if err := query.Order("last_seen DESC").Find(&deviceStatus).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to retrieve device status",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Device status retrieved successfully",
		"data":    deviceStatus,
		"count":   len(deviceStatus),
	})
}

// SendDeviceCommand sends a command to a specific device via MQTT
func (sc *SensorController) SendDeviceCommand(c *gin.Context) {
	deviceID := c.Param("device_id")

	var request requests.DeviceCommandRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Invalid request format",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusBadRequest)
		return
	}

	// Get MQTT service instance
	mqttService := services.GetMQTTService()
	if mqttService == nil || !mqttService.IsConnected() {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "MQTT service not available",
			Errors:  map[string]string{"error": "MQTT service is not connected"},
		}, http.StatusServiceUnavailable)
		return
	}

	// Send command via MQTT
	if err := mqttService.PublishCommand(deviceID, request.Command, request.Payload); err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to send command",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   "Command sent successfully",
		"device_id": deviceID,
		"command":   request.Command,
		"payload":   request.Payload,
		"sent_at":   time.Now(),
	})
}

// GetSensorAlerts gets sensor alerts with filtering options
func (sc *SensorController) GetSensorAlerts(c *gin.Context) {
	var alerts []models.SensorAlert

	query := facades.DB.Model(&models.SensorAlert{})

	// Apply filters
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if kebunName := c.Query("kebun_name"); kebunName != "" {
		query = query.Where("kebun_name = ?", kebunName)
	}

	if severity := c.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	if isResolved := c.Query("is_resolved"); isResolved != "" {
		if resolved, err := strconv.ParseBool(isResolved); err == nil {
			query = query.Where("is_resolved = ?", resolved)
		}
	}

	// Pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	query = query.Limit(limit).Offset(offset).Order("created_at DESC")

	if err := query.Find(&alerts).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to retrieve alerts",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sensor alerts retrieved successfully",
		"data":    alerts,
		"count":   len(alerts),
		"limit":   limit,
		"offset":  offset,
	})
}

// ResolveAlert marks an alert as resolved
func (sc *SensorController) ResolveAlert(c *gin.Context) {
	alertID := c.Param("id")

	var alert models.SensorAlert
	if err := facades.DB.First(&alert, alertID).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Alert not found",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusNotFound)
		return
	}

	now := time.Now()
	if err := facades.DB.Model(&alert).Updates(models.SensorAlert{
		IsResolved: true,
		ResolvedAt: &now,
	}).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to resolve alert",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Alert resolved successfully",
		"data":    alert,
	})
}

// GetSensorStatistics gets aggregated sensor statistics
func (sc *SensorController) GetSensorStatistics(c *gin.Context) {
	query := facades.DB.Model(&models.SensorData{})

	// Apply filters
	if deviceID := c.Query("device_id"); deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	if kebunName := c.Query("kebun_name"); kebunName != "" {
		query = query.Where("kebun_name = ?", kebunName)
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("timestamp >= ?", parsedDate)
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			query = query.Where("timestamp <= ?", parsedDate.Add(24*time.Hour))
		}
	}

	type SensorStats struct {
		AvgNitrogen    float64 `json:"avg_nitrogen"`
		AvgPhosphorus  float64 `json:"avg_phosphorus"`
		AvgPotassium   float64 `json:"avg_potassium"`
		AvgTemperature float64 `json:"avg_temperature"`
		AvgHumidity    float64 `json:"avg_humidity"`
		AvgPH          float64 `json:"avg_ph"`
		MinNitrogen    float64 `json:"min_nitrogen"`
		MaxNitrogen    float64 `json:"max_nitrogen"`
		MinPhosphorus  float64 `json:"min_phosphorus"`
		MaxPhosphorus  float64 `json:"max_phosphorus"`
		MinPotassium   float64 `json:"min_potassium"`
		MaxPotassium   float64 `json:"max_potassium"`
		TotalReadings  int64   `json:"total_readings"`
	}

	var stats SensorStats
	if err := query.Select(`
		AVG(nitrogen) as avg_nitrogen,
		AVG(phosphorus) as avg_phosphorus,
		AVG(potassium) as avg_potassium,
		AVG(temperature) as avg_temperature,
		AVG(humidity) as avg_humidity,
		AVG(ph) as avg_ph,
		MIN(nitrogen) as min_nitrogen,
		MAX(nitrogen) as max_nitrogen,
		MIN(phosphorus) as min_phosphorus,
		MAX(phosphorus) as max_phosphorus,
		MIN(potassium) as min_potassium,
		MAX(potassium) as max_potassium,
		COUNT(*) as total_readings
	`).Scan(&stats).Error; err != nil {
		helpers.ResponseError(c, &helpers.ResponseParams[any]{
			Message: "Failed to retrieve statistics",
			Errors:  map[string]string{"error": err.Error()},
		}, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Sensor statistics retrieved successfully",
		"data":    stats,
	})
}
