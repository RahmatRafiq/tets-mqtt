#!/bin/bash

# Script untuk monitoring MQTT messages
# Usage: ./mqtt_monitor.sh

echo "Monitoring MQTT messages..."
echo "Press Ctrl+C to stop"
echo "========================"

# Monitor semua topic sugar_vestrack
mosquitto_sub -h localhost -t "sugar_vestrack/+/+/+" -v
