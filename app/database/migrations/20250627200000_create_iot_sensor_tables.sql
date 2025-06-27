-- Create sensor_data table for NPK sensor readings
CREATE TABLE sensor_data (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    kebun_name VARCHAR(100) NOT NULL,
    nitrogen DECIMAL(8,2) DEFAULT 0,
    phosphorus DECIMAL(8,2) DEFAULT 0,
    potassium DECIMAL(8,2) DEFAULT 0,
    temperature DECIMAL(5,2) DEFAULT 0,
    humidity DECIMAL(5,2) DEFAULT 0,
    ph DECIMAL(4,2) DEFAULT 0,
    latitude DECIMAL(10,8) DEFAULT 0,
    longitude DECIMAL(11,8) DEFAULT 0,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_device_id (device_id),
    INDEX idx_kebun_name (kebun_name),
    INDEX idx_timestamp (timestamp),
    INDEX idx_deleted_at (deleted_at)
);

-- Create device_status table for tracking device online/offline status
CREATE TABLE device_status (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL UNIQUE,
    kebun_name VARCHAR(100) NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    battery_level INT DEFAULT 0,
    signal_strength INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_device_id (device_id),
    INDEX idx_kebun_name (kebun_name),
    INDEX idx_is_online (is_online)
);

-- Create sensor_alerts table for storing sensor-based alerts
CREATE TABLE sensor_alerts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    kebun_name VARCHAR(100) NOT NULL,
    alert_type VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    severity ENUM('low', 'medium', 'high', 'critical') NOT NULL,
    is_resolved BOOLEAN DEFAULT FALSE,
    resolved_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_device_id (device_id),
    INDEX idx_kebun_name (kebun_name),
    INDEX idx_alert_type (alert_type),
    INDEX idx_severity (severity),
    INDEX idx_is_resolved (is_resolved),
    INDEX idx_created_at (created_at)
);
