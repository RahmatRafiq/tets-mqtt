# 🔧 Developer Guide - IoT Sugar Cane Monitoring System

## Table of Contents
- [Development Setup](#development-setup)
- [Code Architecture](#code-architecture)
- [Database Schema](#database-schema)
- [MQTT Implementation](#mqtt-implementation)
- [API Development Guide](#api-development-guide)
- [Testing Strategy](#testing-strategy)
- [Performance Optimization](#performance-optimization)
- [Security Implementation](#security-implementation)
- [Troubleshooting](#troubleshooting)

---

## Development Setup

### Prerequisites
```bash
# Required tools
- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose
- Git
- Make (optional)
```

### Quick Start
```bash
# 1. Clone repository
git clone <your-repo>
cd tets-mqtt

# 2. Install dependencies
go mod download

# 3. Setup environment
cp .env.example .env

# 4. Start services
docker-compose up -d postgres mqtt

# 5. Run migrations
go run cmd/migrate.go

# 6. Start development server
go run main.go

# 7. Generate Swagger docs
swag init
```

### Development Environment Variables
```env
# Development specific
GO_ENV=development
DEBUG=true
LOG_LEVEL=debug

# Hot reload (with Air)
AIR_BUILD_CMD="go build -o ./tmp/main ."
AIR_BUILD_INCLUDE_EXT="go,tpl,tmpl,html"
```

---

## Code Architecture

### Project Structure Explained
```
tets-mqtt/
├── app/                    # Application core
│   ├── controllers/        # HTTP request handlers
│   ├── middleware/         # HTTP middleware (auth, logging)
│   ├── models/            # Database models (Gorm)
│   ├── requests/          # Request validation structs
│   ├── responses/         # Response formatting
│   ├── services/          # Business logic layer
│   └── helpers/           # Utility functions
├── bootstrap/             # Application initialization
├── cmd/                   # CLI commands (migrate, seed)
├── routes/                # Route definitions
├── docs/                  # Swagger documentation
└── tools/                 # Development tools
```

### Dependency Injection Pattern
```go
// bootstrap/main.go
type AppContainer struct {
    DB          *gorm.DB
    MQTTService *services.MQTTService
    AlertService *services.SensorAlertService
}

func NewAppContainer() *AppContainer {
    db := database.Connection()
    mqttService := services.NewMQTTService(db)
    alertService := services.NewSensorAlertService(db)
    
    return &AppContainer{
        DB:           db,
        MQTTService:  mqttService,
        AlertService: alertService,
    }
}
```

### Service Layer Pattern
```go
// app/services/sensor_service.go
type SensorService struct {
    db           *gorm.DB
    alertService *SensorAlertService
}

func (s *SensorService) ProcessSensorData(data *models.SensorData) error {
    // 1. Validate data
    if err := s.validateSensorData(data); err != nil {
        return err
    }
    
    // 2. Save to database
    if err := s.db.Create(data).Error; err != nil {
        return err
    }
    
    // 3. Check for alerts
    go s.alertService.CheckThresholds(data)
    
    return nil
}
```

---

## Database Schema

### Migration System
```go
// cmd/migrate.go
func runMigrations() {
    files, _ := filepath.Glob("app/database/migrations/*.sql")
    sort.Strings(files)
    
    for _, file := range files {
        content, _ := ioutil.ReadFile(file)
        db.Exec(string(content))
        log.Printf("Applied migration: %s", file)
    }
}
```

### Model Relationships
```go
// app/models/sensor_data.go
type SensorData struct {
    ID             uint      `json:"id" gorm:"primaryKey"`
    DeviceID       string    `json:"device_id" gorm:"index;not null"`
    FarmName       string    `json:"farm_name" gorm:"not null"`
    Location       string    `json:"location"`
    Temperature    float64   `json:"temperature"`
    Humidity       float64   `json:"humidity"`
    SoilPH         float64   `json:"soil_ph"`
    LightIntensity float64   `json:"light_intensity"`
    Timestamp      time.Time `json:"timestamp"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
    
    // Relationships
    Device *DeviceStatus `json:"device,omitempty" gorm:"foreignKey:DeviceID;references:DeviceID"`
    Alerts []SensorAlert `json:"alerts,omitempty" gorm:"foreignKey:DeviceID;references:DeviceID"`
}
```

### Advanced Queries
```go
// Complex queries with preloading
func (s *SensorService) GetSensorDataWithAlerts(deviceID string) ([]models.SensorData, error) {
    var data []models.SensorData
    
    err := s.db.
        Preload("Device").
        Preload("Alerts", "is_resolved = ?", false).
        Where("device_id = ?", deviceID).
        Order("timestamp DESC").
        Limit(100).
        Find(&data).Error
        
    return data, err
}

// Aggregation queries
func (s *SensorService) GetDailyAverages(farmName string, days int) ([]DailyAverage, error) {
    var results []DailyAverage
    
    err := s.db.Raw(`
        SELECT 
            DATE(timestamp) as date,
            AVG(temperature) as avg_temperature,
            AVG(humidity) as avg_humidity,
            AVG(soil_ph) as avg_soil_ph,
            AVG(light_intensity) as avg_light_intensity
        FROM sensor_data 
        WHERE farm_name = ? 
        AND timestamp >= NOW() - INTERVAL ? DAY
        GROUP BY DATE(timestamp)
        ORDER BY date DESC
    `, farmName, days).Scan(&results).Error
    
    return results, err
}
```

---

## MQTT Implementation

### MQTT Client Configuration
```go
// app/services/mqtt_service.go
func NewMQTTService(db *gorm.DB) *MQTTService {
    opts := mqtt.NewClientOptions()
    opts.AddBroker(os.Getenv("MQTT_BROKER"))
    opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
    opts.SetCleanSession(true)
    opts.SetOrderMatters(false)
    opts.SetKeepAlive(30 * time.Second)
    opts.SetPingTimeout(10 * time.Second)
    opts.SetConnectTimeout(10 * time.Second)
    opts.SetMaxReconnectInterval(10 * time.Second)
    opts.SetAutoReconnect(true)
    
    // Connection handlers
    opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
        log.Printf("MQTT connection lost: %v", err)
    })
    
    opts.SetReconnectingHandler(func(client mqtt.Client, opts *mqtt.ClientOptions) {
        log.Println("MQTT attempting to reconnect...")
    })
    
    client := mqtt.NewClient(opts)
    return &MQTTService{client: client, db: db}
}
```

### Message Processing Pipeline
```go
func (m *MQTTService) processSensorData(topic string, payload []byte) {
    // 1. Parse device ID from topic
    deviceID := m.extractDeviceID(topic)
    
    // 2. Unmarshal JSON payload
    var sensorPayload SensorDataPayload
    if err := json.Unmarshal(payload, &sensorPayload); err != nil {
        log.Printf("Invalid sensor payload: %v", err)
        return
    }
    
    // 3. Convert to model
    sensorData := sensorPayload.ToModel()
    
    // 4. Validate data
    if err := m.validateSensorData(sensorData); err != nil {
        log.Printf("Invalid sensor data: %v", err)
        return
    }
    
    // 5. Save to database (with retry)
    if err := m.saveSensorDataWithRetry(sensorData, 3); err != nil {
        log.Printf("Failed to save sensor data: %v", err)
        return
    }
    
    // 6. Process alerts asynchronously
    go m.alertService.CheckThresholds(sensorData)
}
```

### Error Handling and Retry Logic
```go
func (m *MQTTService) saveSensorDataWithRetry(data *models.SensorData, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        if err := m.db.Create(data).Error; err != nil {
            lastErr = err
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        return nil
    }
    
    return fmt.Errorf("failed after %d retries: %v", maxRetries, lastErr)
}
```

---

## API Development Guide

### Request Validation
```go
// app/requests/sensor_request.go
type CreateSensorDataRequest struct {
    DeviceID       string  `json:"device_id" binding:"required,min=3,max=50"`
    FarmName       string  `json:"farm_name" binding:"required,min=3,max=100"`
    Location       string  `json:"location" binding:"max=100"`
    Temperature    float64 `json:"temperature" binding:"min=-50,max=100"`
    Humidity       float64 `json:"humidity" binding:"min=0,max=100"`
    SoilPH         float64 `json:"soil_ph" binding:"min=0,max=14"`
    LightIntensity float64 `json:"light_intensity" binding:"min=0"`
    Timestamp      string  `json:"timestamp" binding:"required"`
}

func (r *CreateSensorDataRequest) Validate() error {
    // Custom validation logic
    if r.Temperature < -50 || r.Temperature > 100 {
        return errors.New("temperature must be between -50 and 100")
    }
    
    if _, err := time.Parse(time.RFC3339, r.Timestamp); err != nil {
        return errors.New("invalid timestamp format")
    }
    
    return nil
}
```

### Response Formatting
```go
// app/handlers/response_handler.go
type APIResponse struct {
    Status     string      `json:"status"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    Pagination *Pagination `json:"pagination,omitempty"`
    Error      *APIError   `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, APIResponse{
        Status:  "success",
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
    response := APIResponse{
        Status:  "error",
        Message: message,
    }
    
    if err != nil {
        response.Error = &APIError{
            Code:    statusCode,
            Details: err.Error(),
        }
    }
    
    c.JSON(statusCode, response)
}
```

### Middleware Implementation
```go
// app/middleware/logger_middleware.go
func LoggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}

// app/middleware/cors_middleware.go
func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    })
}
```

---

## Testing Strategy

### Unit Tests
```go
// app/services/sensor_service_test.go
func TestSensorService_ProcessSensorData(t *testing.T) {
    // Setup
    db := setupTestDB()
    alertService := NewMockAlertService()
    service := NewSensorService(db, alertService)
    
    // Test data
    sensorData := &models.SensorData{
        DeviceID:       "TEST_001",
        FarmName:       "Test Farm",
        Temperature:    25.5,
        Humidity:       60.0,
        SoilPH:         6.5,
        LightIntensity: 800,
        Timestamp:      time.Now(),
    }
    
    // Execute
    err := service.ProcessSensorData(sensorData)
    
    // Assert
    assert.NoError(t, err)
    assert.True(t, alertService.CheckThresholdsCalled)
    
    // Verify database
    var saved models.SensorData
    db.First(&saved, "device_id = ?", "TEST_001")
    assert.Equal(t, sensorData.Temperature, saved.Temperature)
}
```

### Integration Tests
```go
// app/controllers/sensor_controller_test.go
func TestSensorController_GetSensorData(t *testing.T) {
    // Setup test server
    router := setupTestRouter()
    w := httptest.NewRecorder()
    
    // Create test data
    seedTestData()
    
    // Execute request
    req, _ := http.NewRequest("GET", "/api/v1/sensors/data?device_id=TEST_001", nil)
    router.ServeHTTP(w, req)
    
    // Assert response
    assert.Equal(t, 200, w.Code)
    
    var response APIResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, "success", response.Status)
    assert.NotNil(t, response.Data)
}
```

### MQTT Testing
```go
// Test MQTT message processing
func TestMQTTService_ProcessMessage(t *testing.T) {
    service := setupTestMQTTService()
    
    payload := `{
        "device_id": "TEST_001",
        "farm_name": "Test Farm",
        "temperature": 25.5,
        "humidity": 60.0,
        "soil_ph": 6.5,
        "light_intensity": 800,
        "timestamp": "2025-06-27T10:30:00Z"
    }`
    
    err := service.ProcessSensorData("sensor/TEST_001/data", []byte(payload))
    assert.NoError(t, err)
    
    // Verify data saved
    var saved models.SensorData
    service.db.First(&saved, "device_id = ?", "TEST_001")
    assert.Equal(t, 25.5, saved.Temperature)
}
```

---

## Performance Optimization

### Database Optimization
```go
// Connection pooling
func setupDatabase() *gorm.DB {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db
}

// Batch insert for high-volume data
func (s *SensorService) BatchInsertSensorData(data []models.SensorData) error {
    batchSize := 1000
    
    for i := 0; i < len(data); i += batchSize {
        end := i + batchSize
        if end > len(data) {
            end = len(data)
        }
        
        if err := s.db.CreateInBatches(data[i:end], batchSize).Error; err != nil {
            return err
        }
    }
    
    return nil
}
```

### Caching Strategy
```go
// Redis caching for frequently accessed data
type CacheService struct {
    redis *redis.Client
}

func (c *CacheService) GetSensorStats(deviceID string) (*SensorStats, error) {
    key := fmt.Sprintf("sensor_stats:%s", deviceID)
    
    // Try cache first
    cached, err := c.redis.Get(key).Result()
    if err == nil {
        var stats SensorStats
        json.Unmarshal([]byte(cached), &stats)
        return &stats, nil
    }
    
    // Cache miss - fetch from database
    stats, err := c.calculateSensorStats(deviceID)
    if err != nil {
        return nil, err
    }
    
    // Cache result for 5 minutes
    data, _ := json.Marshal(stats)
    c.redis.Set(key, data, 5*time.Minute)
    
    return stats, nil
}
```

### Goroutine Management
```go
// Worker pool for processing MQTT messages
type WorkerPool struct {
    workers   int
    taskQueue chan Task
    wg        sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:   workers,
        taskQueue: make(chan Task, workers*2),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    
    for task := range wp.taskQueue {
        task.Execute()
    }
}
```

---

## Security Implementation

### JWT Authentication
```go
// app/middleware/auth_middleware.go
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        
        // Remove "Bearer " prefix
        if strings.HasPrefix(token, "Bearer ") {
            token = token[7:]
        }
        
        claims, err := validateJWT(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

### Rate Limiting
```go
// Rate limiting middleware
func RateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(time.Minute), 100) // 100 requests per minute
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": "60s",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### Input Validation and Sanitization
```go
// SQL injection prevention
func (s *SensorService) GetSensorData(filters map[string]interface{}) ([]models.SensorData, error) {
    query := s.db.Model(&models.SensorData{})
    
    // Safe parameter binding
    if deviceID, ok := filters["device_id"]; ok {
        query = query.Where("device_id = ?", deviceID)
    }
    
    if farmName, ok := filters["farm_name"]; ok {
        query = query.Where("farm_name = ?", farmName)
    }
    
    // Safe date range filtering
    if startDate, ok := filters["start_date"]; ok {
        query = query.Where("timestamp >= ?", startDate)
    }
    
    var data []models.SensorData
    err := query.Find(&data).Error
    return data, err
}
```

---

## Troubleshooting

### Common Issues and Solutions

#### 1. MQTT Connection Issues
```go
// Debug MQTT connectivity
func debugMQTTConnection() {
    client := mqtt.NewClient(opts)
    
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        log.Printf("MQTT connection failed: %v", token.Error())
        
        // Check broker availability
        conn, err := net.Dial("tcp", brokerHost)
        if err != nil {
            log.Printf("Broker unreachable: %v", err)
        } else {
            conn.Close()
            log.Println("Broker reachable, check credentials")
        }
    }
}
```

#### 2. Database Performance Issues
```sql
-- Check slow queries
SELECT query, mean_time, calls 
FROM pg_stat_statements 
WHERE mean_time > 1000 
ORDER BY mean_time DESC;

-- Add missing indexes
CREATE INDEX CONCURRENTLY idx_sensor_data_device_timestamp 
ON sensor_data(device_id, timestamp DESC);

-- Check connection pool
SELECT 
    numbackends,
    datname 
FROM pg_stat_database 
WHERE datname = 'sugar_cane_iot';
```

#### 3. Memory Leaks
```go
// Monitor goroutines
func monitorGoroutines() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            count := runtime.NumGoroutine()
            log.Printf("Active goroutines: %d", count)
            
            if count > 1000 {
                log.Println("WARNING: High goroutine count detected")
                // Trigger graceful shutdown or investigation
            }
        }
    }
}
```

### Logging and Monitoring
```go
// Structured logging
func setupLogger() *logrus.Logger {
    logger := logrus.New()
    
    logger.SetFormatter(&logrus.JSONFormatter{
        TimestampFormat: time.RFC3339,
    })
    
    logger.SetLevel(logrus.InfoLevel)
    
    // Add hooks for external monitoring
    hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_INFO, "")
    if err == nil {
        logger.Hooks.Add(hook)
    }
    
    return logger
}

// Application metrics
type Metrics struct {
    MQTTMessagesReceived prometheus.Counter
    APIRequestDuration   prometheus.Histogram
    DatabaseConnections  prometheus.Gauge
}

func setupMetrics() *Metrics {
    return &Metrics{
        MQTTMessagesReceived: prometheus.NewCounter(prometheus.CounterOpts{
            Name: "mqtt_messages_received_total",
            Help: "Total number of MQTT messages received",
        }),
        APIRequestDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
            Name: "api_request_duration_seconds",
            Help: "Duration of API requests",
        }),
        DatabaseConnections: prometheus.NewGauge(prometheus.GaugeOpts{
            Name: "database_connections_active",
            Help: "Number of active database connections",
        }),
    }
}
```

---

## Best Practices

### Code Organization
- Keep controllers thin - move business logic to services
- Use dependency injection for testability
- Implement proper error handling and logging
- Follow Go naming conventions and project structure

### Database
- Use migrations for schema changes
- Implement proper indexing strategy
- Use connection pooling
- Implement batch operations for high-volume data

### MQTT
- Implement proper error handling and reconnection logic
- Use appropriate QoS levels
- Implement message deduplication
- Monitor broker health and performance

### API Design
- Follow RESTful conventions
- Implement proper pagination
- Use appropriate HTTP status codes
- Implement request validation and sanitization

### Testing
- Write unit tests for all business logic
- Implement integration tests for API endpoints
- Use test databases for testing
- Mock external dependencies

---

**Created by**: Development Team  
**Date**: June 27, 2025  
**Version**: 1.0.0
