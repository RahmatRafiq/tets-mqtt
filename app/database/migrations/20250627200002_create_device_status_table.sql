-- Create device_status table for tracking device online/offline status
CREATE TABLE device_status (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    device_id VARCHAR(100) NOT NULL UNIQUE,
    farm_name VARCHAR(100) NOT NULL,
    is_online BOOLEAN DEFAULT FALSE,
    last_seen TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    battery_level DECIMAL(5,2) DEFAULT 0,
    signal_strength INT DEFAULT 0,
    firmware_version VARCHAR(50) DEFAULT '',
    location VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_device_id (device_id),
    INDEX idx_farm_name (farm_name),
    INDEX idx_is_online (is_online)
);
