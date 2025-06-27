# Sugar VesTrack IoT MQTT Implementation

## Overview
This implementation provides MQTT connectivity for the Sugar VesTrack IoT monitoring system based on the architecture diagram provided. The system handles NPK sensor data from ESP32-S3 devices deployed across multiple kebun (plantation areas).

## Architecture Components

### Hardware Layer
- **ESP32-S3 Boards**: Primary microcontrollers with built-in WiFi/Bluetooth
- **NPK RS485 Sensors**: Soil nutrient sensors measuring Nitrogen, Phosphorus, and Potassium
- **LoRaWAN Gateway**: DLOS5N-EC25 for wide-area connectivity
- **4G Connectivity**: For internet connection to cloud services

### MQTT Topics Structure

#### Sensor Data Topic
```
sugar_vestrack/sensor/{device_id}/data
```

**Payload Example:**
```json
{
  "device_id": "ESP32_001",
  "kebun_name": "Kebun A",
  "nitrogen": 25.5,
  "phosphorus": 18.2,
  "potassium": 165.0,
  "temperature": 28.5,
  "humidity": 65.2,
  "ph": 6.8,
  "latitude": -7.250445,
  "longitude": 112.768845,
  "battery_level": 85,
  "signal_strength": -65,
  "timestamp": 1640995200
}
```

#### Device Status Topic
```
sugar_vestrack/device/{device_id}/status
```

**Payload Example:**
```json
{
  "device_id": "ESP32_001",
  "kebun_name": "Kebun A",
  "is_online": true,
  "battery_level": 85,
  "signal_strength": -65,
  "timestamp": 1640995200
}
```

#### Device Command Topic
```
sugar_vestrack/device/{device_id}/command
```

**Payload Example:**
```json
{
  "command": "set_interval",
  "payload": {
    "interval_minutes": 30
  },
  "timestamp": 1640995200
}
```

## API Endpoints

### Get Sensor Data
```
GET /api/sensor/data
```

**Query Parameters:**
- `device_id`: Filter by specific device
- `kebun_name`: Filter by kebun name
- `start_date`: Start date (YYYY-MM-DD)
- `end_date`: End date (YYYY-MM-DD)
- `limit`: Number of records (default: 50)
- `offset`: Pagination offset (default: 0)

### Get Latest Sensor Data
```
GET /api/sensor/latest?device_id={device_id}
```

### Get Device Status
```
GET /api/sensor/devices/status?device_id={device_id}
```

### Send Device Command
```
POST /api/sensor/devices/{device_id}/command
```

**Request Body:**
```json
{
  "command": "restart",
  "payload": {}
}
```

## Available Commands

### Device Commands
- `restart`: Restart the ESP32 device
- `calibrate`: Calibrate NPK sensors
- `set_interval`: Set data transmission interval
- `sleep`: Put device into sleep mode
- `wake_up`: Wake device from sleep mode

### Configuration Commands
- `set_wifi`: Update WiFi credentials
- `set_mqtt`: Update MQTT broker settings
- `set_threshold`: Set alert thresholds for NPK values

## Database Tables

### sensor_data
Stores all NPK sensor readings with timestamps and location data.

### device_status
Tracks online/offline status, battery levels, and signal strength.

### sensor_alerts
Stores automatically generated alerts based on sensor thresholds.

## Alert System

The system automatically generates alerts for:
- **Low Nitrogen**: < 20 mg/kg
- **Low Phosphorus**: < 15 mg/kg
- **Low Potassium**: < 150 mg/kg
- **pH Abnormal**: < 6.0 or > 8.0 (optimal for sugar cane: 6.0-8.0)
- **Device Offline**: No data received for > 1 hour
- **Low Battery**: < 20%

## Setup Instructions

### 1. Environment Configuration
Configure MQTT settings in your `.env` file:

```bash
# MQTT Configuration for Sugar VesTrack IoT System
MQTT_BROKER_URL=tcp://broker.hivemq.com:1883
MQTT_CLIENT_ID=sugar_vestrack_server
MQTT_USERNAME=your_username
MQTT_PASSWORD=your_password
```

### 2. Database Migration
Run the migration to create sensor tables:

```bash
go run main.go migrate
```

### 3. Start the Server
```bash
go run main.go
```

The MQTT service will automatically connect and start listening for sensor data.

## ESP32 Integration

### Sample ESP32 Code Structure
```cpp
#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

// NPK Sensor reading
void readNPKSensor() {
  // Read from RS485 NPK sensor
  // Format and publish to MQTT
}

// Publish sensor data
void publishSensorData() {
  StaticJsonDocument<512> doc;
  doc["device_id"] = "ESP32_001";
  doc["kebun_name"] = "Kebun A";
  doc["nitrogen"] = nitrogen_value;
  doc["phosphorus"] = phosphorus_value;
  doc["potassium"] = potassium_value;
  // ... other sensor values
  
  String payload;
  serializeJson(doc, payload);
  client.publish("sugar_vestrack/sensor/ESP32_001/data", payload.c_str());
}
```

## Monitoring and Troubleshooting

### MQTT Connection Status
Check the server logs for MQTT connection status:
```
MQTT client connected to broker
Successfully subscribed to topic: sugar_vestrack/sensor/+/data
```

### Test MQTT Connection
Use MQTT client tools to test:
```bash
# Subscribe to all sensor data
mosquitto_sub -h broker.hivemq.com -t "sugar_vestrack/sensor/+/data"

# Publish test data
mosquitto_pub -h broker.hivemq.com -t "sugar_vestrack/sensor/TEST001/data" -m '{"device_id":"TEST001","nitrogen":25.5}'
```

## Security Considerations

1. **MQTT Authentication**: Use username/password authentication
2. **TLS/SSL**: Use secure connections (ssl://) for production
3. **Topic ACL**: Implement topic-based access control
4. **Device Certificates**: Use client certificates for device authentication
5. **Data Encryption**: Encrypt sensitive payload data

## Performance Optimization

1. **QoS Levels**: Use QoS 1 for sensor data to ensure delivery
2. **Retained Messages**: Use retained messages for device status
3. **Connection Pooling**: Implement connection pooling for high-volume scenarios
4. **Database Indexing**: Ensure proper indexing on frequently queried fields

## Future Enhancements

1. **Real-time Dashboard**: WebSocket integration for live data updates
2. **Data Analytics**: Historical data analysis and trends
3. **Machine Learning**: Predictive analytics for crop health
4. **Mobile App**: React Native app for field workers
5. **Report Generation**: Automated PDF reports for farm management
