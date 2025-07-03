# MQTT Testing Guide

## 🚀 Cara Menjalankan Semua Komponen

### 1. Menjalankan MQTT Broker (Mosquitto)
```bash
# Cek status Mosquitto
sudo systemctl status mosquitto

# Jika belum berjalan, start Mosquitto
sudo systemctl start mosquitto

# Untuk auto-start saat boot
sudo systemctl enable mosquitto

# Restart jika diperlukan
sudo systemctl restart mosquitto
```

### 2. Menjalankan Database MariaDB
```bash
# Cek status MariaDB
sudo systemctl status mariadb

# Jika belum berjalan, start MariaDB
sudo systemctl start mariadb

# Untuk auto-start saat boot
sudo systemctl enable mariadb

# Test koneksi database
mysql -u rahmat -p -h localhost -e "SHOW DATABASES;"
```

### 3. Menjalankan Aplikasi Golang
```bash
# Masuk ke directory project
cd /home/rahmat/golang-project/tets-mqtt

# Jalankan aplikasi
go run main.go

# Atau build dulu, lalu jalankan
go build -o sugar_vestrack main.go
./sugar_vestrack
```

### 4. Cek Status Semua Komponen
```bash
# Cek semua service sekaligus
sudo systemctl status mosquitto mariadb

# Cek port yang digunakan
ss -tlnp | grep -E "(1883|3306|8080)"
```

## 🔧 Langkah-langkah Testing

### Step 1: Pastikan Semua Service Berjalan
```bash
# Terminal 1: Cek status service
sudo systemctl status mosquitto mariadb

# Jika ada yang tidak berjalan, start service tersebut
sudo systemctl start mosquitto mariadb
```

### Step 2: Jalankan Aplikasi Golang
```bash
# Terminal 2: Jalankan aplikasi
cd /home/rahmat/golang-project/tets-mqtt
go run main.go
```

### Step 3: Monitoring MQTT Messages (Opsional)
```bash
# Terminal 3: Monitor MQTT messages
./mqtt_monitor.sh

# Atau manual:
mosquitto_sub -h localhost -t "sugar_vestrack/+/+/+" -v
```

### Step 4: Kirim Data Test
```bash
# Terminal 4: Kirim data sensor
./mqtt_test.sh device001 28.5 65.3 42.7

# Atau manual:
mosquitto_pub -h localhost -t "sugar_vestrack/sensor/device001/data" -m '{"deviceId":"device001","temperature":28.5,"humidity":65.3,"soilMoisture":42.7,"timestamp":1720027620}'
```

### Step 5: Cek Data via API
```bash
# Cek data sensor yang masuk
curl -X GET "http://localhost:8080/api/sensor/data" -H "Content-Type: application/json"
```

## 🛠️ Quick Start Commands

### Menjalankan Semuanya dalam 3 Command:
```bash
# 1. Start semua service yang diperlukan
sudo systemctl start mosquitto mariadb

# 2. Masuk ke project directory dan jalankan aplikasi
cd /home/rahmat/golang-project/tets-mqtt && go run main.go

# 3. Test kirim data sensor (jalankan di terminal lain)
./mqtt_test.sh device001 25.5 60.2 45.8
```

### One-liner untuk Cek Status:
```bash
# Cek semua status sekaligus
echo "=== MQTT Broker ===" && sudo systemctl is-active mosquitto && echo "=== Database ===" && sudo systemctl is-active mariadb && echo "=== Ports ===" && ss -tlnp | grep -E "(1883|3306|8080)"
```

## 📱 Command Cheat Sheet

### MQTT Commands:
```bash
# Publish pesan
mosquitto_pub -h localhost -t "topic/name" -m "message"

# Subscribe ke topic
mosquitto_sub -h localhost -t "topic/name" -v

# Monitor semua topic sugar_vestrack
mosquitto_sub -h localhost -t "sugar_vestrack/+/+/+" -v
```

### API Testing:
```bash
# Get semua data sensor
curl -X GET "http://localhost:8080/api/sensor/data"

# Get data sensor terbaru
curl -X GET "http://localhost:8080/api/sensor/latest"

# Get statistik sensor
curl -X GET "http://localhost:8080/api/sensor/statistics"

# Health check
curl -X GET "http://localhost:8080/health"
```

### Database Commands:
```bash
# Login ke database
mysql -u rahmat -p

# Cek data sensor di database
mysql -u rahmat -p -e "USE golang_starter_kit_2025; SELECT * FROM sensor_data ORDER BY created_at DESC LIMIT 5;"
```

## 🚨 Troubleshooting

### Jika MQTT tidak bisa connect:
```bash
# Restart Mosquitto
sudo systemctl restart mosquitto

# Cek konfigurasi Mosquitto
sudo nano /etc/mosquitto/mosquitto.conf

# Test koneksi manual
mosquitto_pub -h localhost -t test -m "hello"
```

### Jika Database tidak bisa connect:
```bash
# Restart MariaDB
sudo systemctl restart mariadb

# Reset password jika perlu
sudo mysql_secure_installation

# Test koneksi
mysql -u rahmat -p -e "SELECT 1;"
```

### Jika Aplikasi error:
```bash
# Cek dependencies
go mod tidy

# Build ulang
go clean && go build

# Cek log error
go run main.go 2>&1 | tee app.log
```

## 🔗 Tanda-tanda Koneksi MQTT Berhasil

### 1. Log Output saat Aplikasi Start
Ketika aplikasi Golang berhasil terhubung dengan MQTT, Anda akan melihat log berikut:

```
2025/07/03 22:06:51 Successfully connected to MQTT broker
2025/07/03 22:06:51 MQTT service initialized successfully
2025/07/03 22:06:51 MQTT client connected to broker
2025/07/03 22:06:51 Successfully subscribed to topic: sugar_vestrack/sensor/+/data
2025/07/03 22:06:51 Successfully subscribed to topic: sugar_vestrack/device/+/status
2025/07/03 22:06:51 Successfully subscribed to topic: sugar_vestrack/device/+/alert
```

### 2. Tidak Ada Error Message
Jika koneksi gagal, Anda akan melihat error seperti:
```
2025/07/03 22:06:51 Failed to connect to MQTT broker: network Error: dial tcp 127.0.0.1:1883: connect: connection refused
2025/07/03 22:06:51 MQTT functionality will be disabled
```

### 3. Test Koneksi Real-time
Untuk memastikan koneksi benar-benar berfungsi:

```bash
# Terminal 1: Jalankan aplikasi Golang
go run main.go

# Terminal 2: Kirim pesan test
mosquitto_pub -h localhost -t "sugar_vestrack/sensor/test/data" -m '{"deviceId":"test","temperature":25.0,"humidity":60.0,"soilMoisture":45.0,"timestamp":1720027620}'

# Jika berhasil, di Terminal 1 akan muncul:
# 2025/07/03 22:13:08 Received sensor data from topic: sugar_vestrack/sensor/test/data
```

### 4. Cek Status Client MQTT
Aplikasi akan menampilkan status koneksi MQTT:

```bash
# Status berhasil:
✅ "MQTT client connected to broker"
✅ "Successfully subscribed to topic: ..."

# Status gagal:
❌ "Failed to connect to MQTT broker"
❌ "MQTT functionality will be disabled"
```

### 5. Cek Data Masuk ke Database
Setelah mengirim data via MQTT, cek apakah data tersimpan:

```bash
# Test API endpoint
curl -X GET "http://localhost:8080/api/sensor/data" -H "Content-Type: application/json"

# Atau cek langsung di database
mysql -u rahmat -p -e "USE golang_starter_kit_2025; SELECT * FROM sensor_data ORDER BY created_at DESC LIMIT 3;"
```

### 6. Monitoring Real-time
Untuk monitoring koneksi secara real-time:

```bash
# Terminal monitoring - akan menampilkan semua pesan MQTT
mosquitto_sub -h localhost -t "sugar_vestrack/+/+/+" -v

# Jika ada pesan masuk, berarti koneksi aktif
```

### 7. Health Check API
Cek status aplikasi via API:

```bash
curl -X GET "http://localhost:8080/health"

# Response berhasil:
{"status":"ok","timestamp":"2025-07-03T22:06:51Z"}
```

### 8. Indikator di Code
Dalam kode `bootstrap/mqtt.go`, fungsi `Connect()` akan return:
- `nil` jika berhasil → Log: "MQTT service initialized successfully"
- `error` jika gagal → Log: "Failed to connect to MQTT broker"

### 9. Process Monitoring
Cek apakah aplikasi subscribe ke topic dengan benar:

```bash
# Cek proses MQTT client
ps aux | grep -i mqtt

# Cek koneksi network
ss -tlnp | grep 1883
```

## 🎉 Status Koneksi
- ✅ Mosquitto MQTT Broker: Active dan listening di port 1883
- ✅ MariaDB Database: Terhubung dengan user rahmat
- ✅ Golang Application: Berjalan di port 8080
- ✅ MQTT Subscription: Berhasil subscribe ke 3 topik utama
- ✅ Data Persistence: Data sensor berhasil disimpan ke database
