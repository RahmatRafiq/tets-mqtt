CREATE TABLE IF NOT EXISTS sensor_alerts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key for sensor alert records',
    device_id VARCHAR(100) NOT NULL COMMENT 'Device identifier that triggered the alert',
    farm_name VARCHAR(100) NOT NULL COMMENT 'Name of the farm where the alert occurred',
    alert_type VARCHAR(50) NOT NULL COMMENT 'Type of alert: nitrogen_low, phosphorus_low, potassium_low, ph_abnormal, etc.',
    message TEXT NOT NULL COMMENT 'Human-readable alert message describing the issue',
    severity ENUM('low', 'medium', 'high', 'critical') NOT NULL COMMENT 'Alert severity level',
    sensor_value DECIMAL(10,2) DEFAULT NULL COMMENT 'Actual sensor value that triggered the alert',
    threshold_value DECIMAL(10,2) DEFAULT NULL COMMENT 'Threshold value that was violated',
    is_resolved BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Whether the alert has been resolved/acknowledged',
    resolved_at TIMESTAMP NULL COMMENT 'Timestamp when the alert was resolved',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Alert creation timestamp',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Alert last update timestamp',
    
    -- Indexes for better query performance
    INDEX idx_sensor_alerts_device_id (device_id) COMMENT 'Index on device_id for device-specific alert queries',
    INDEX idx_sensor_alerts_farm_name (farm_name) COMMENT 'Index on farm_name for farm-specific alert queries',
    INDEX idx_sensor_alerts_alert_type (alert_type) COMMENT 'Index on alert_type for filtering by alert category',
    INDEX idx_sensor_alerts_severity (severity) COMMENT 'Index on severity for filtering by alert severity',
    INDEX idx_sensor_alerts_is_resolved (is_resolved) COMMENT 'Index on is_resolved for filtering resolved/unresolved alerts',
    INDEX idx_sensor_alerts_created_at (created_at) COMMENT 'Index on created_at for time-based alert queries',
    INDEX idx_sensor_alerts_device_resolved (device_id, is_resolved) COMMENT 'Composite index for device-specific resolution status',
    INDEX idx_sensor_alerts_farm_severity (farm_name, severity) COMMENT 'Composite index for farm-specific severity queries',
    INDEX idx_sensor_alerts_type_severity (alert_type, severity) COMMENT 'Composite index for alert type and severity queries'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Table storing sensor-based alerts and threshold violations';
