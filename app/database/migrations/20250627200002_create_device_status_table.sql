CREATE TABLE IF NOT EXISTS device_status (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'Primary key for device status records',
    device_id VARCHAR(100) NOT NULL UNIQUE COMMENT 'Unique identifier for the IoT device',
    farm_name VARCHAR(100) NOT NULL COMMENT 'Name of the farm where the device is located',
    is_online BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Whether the device is currently online and communicating',
    last_seen TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Last time the device was seen/communicated',
    battery_level DECIMAL(5,2) NOT NULL DEFAULT 0.00 COMMENT 'Battery level percentage (0.00-100.00)',
    signal_strength INT NOT NULL DEFAULT 0 COMMENT 'Signal strength in dBm (negative values)',
    firmware_version VARCHAR(50) DEFAULT '' COMMENT 'Current firmware version of the device',
    location VARCHAR(255) DEFAULT '' COMMENT 'Physical location description of the device',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record last update timestamp',
    
    UNIQUE KEY uk_device_status_device_id (device_id) COMMENT 'Unique constraint on device_id',
    INDEX idx_device_status_farm_name (farm_name) COMMENT 'Index on farm_name for filtering by farm',
    INDEX idx_device_status_is_online (is_online) COMMENT 'Index on is_online for filtering online/offline devices',
    INDEX idx_device_status_last_seen (last_seen) COMMENT 'Index on last_seen for time-based queries',
    INDEX idx_device_status_battery_level (battery_level) COMMENT 'Index on battery_level for low battery queries',
    INDEX idx_device_status_farm_online (farm_name, is_online) COMMENT 'Composite index for farm-specific online status queries'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Table tracking operational status of IoT devices';
