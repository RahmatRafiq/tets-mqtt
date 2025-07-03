#!/bin/bash

# Script untuk mengirim data sensor test via MQTT
# Usage: ./mqtt_test.sh [device_id] [temperature] [humidity] [soil_moisture]

DEVICE_ID=${1:-"device001"}
TEMPERATURE=${2:-$(echo "scale=1; 20 + $RANDOM % 20" | bc)}
HUMIDITY=${3:-$(echo "scale=1; 50 + $RANDOM % 30" | bc)}
SOIL_MOISTURE=${4:-$(echo "scale=1; 30 + $RANDOM % 40" | bc)}
TIMESTAMP=$(date +%s)

# JSON payload
PAYLOAD=$(cat <<EOF
{
  "deviceId": "$DEVICE_ID",
  "temperature": $TEMPERATURE,
  "humidity": $HUMIDITY,
  "soilMoisture": $SOIL_MOISTURE,
  "timestamp": $TIMESTAMP
}
EOF
)

echo "Sending MQTT message:"
echo "Topic: sugar_vestrack/sensor/$DEVICE_ID/data"
echo "Payload: $PAYLOAD"

# Kirim pesan MQTT
mosquitto_pub -h localhost -t "sugar_vestrack/sensor/$DEVICE_ID/data" -m "$PAYLOAD"

echo "Message sent successfully!"
