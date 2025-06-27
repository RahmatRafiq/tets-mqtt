-- Create sensor_data table for NPK sensor readings
CREATE TABLE sensor_data (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL,
    farm_name VARCHAR(100) NOT NULL,
    nitrogen DECIMAL(8,2) DEFAULT 0,
    phosphorus DECIMAL(8,2) DEFAULT 0,
    potassium DECIMAL(8,2) DEFAULT 0,
    temperature DECIMAL(5,2) DEFAULT 0,
    humidity DECIMAL(5,2) DEFAULT 0,
    ph DECIMAL(4,2) DEFAULT 0,
    latitude DECIMAL(10,8) DEFAULT 0,
    longitude DECIMAL(11,8) DEFAULT 0,
    location VARCHAR(255) DEFAULT '',
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    INDEX idx_device_id (device_id),
    INDEX idx_farm_name (farm_name),
    INDEX idx_timestamp (timestamp),
    INDEX idx_deleted_at (deleted_at)
);
