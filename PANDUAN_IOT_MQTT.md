# 📖 Panduan Lengkap: Sistem IoT Monitoring Tanaman Tebu dengan MQTT

## 🌟 Pengantar untuk Pemula

Dokumen ini menjelaskan secara sederhana tentang sistem monitoring tanaman tebu yang menggunakan teknologi IoT (Internet of Things) yang telah kita bangun. Sistem ini dapat memantau kondisi tanah secara real-time menggunakan sensor-sensor khusus.

---

## 🤔 Apa itu IoT (Internet of Things)?

**IoT** adalah konsep dimana perangkat-perangkat (seperti sensor, kamera, dll) dapat terhubung ke internet dan saling berkomunikasi. 

**Contoh sederhana:**
- Seperti lampu yang bisa dinyalakan dari smartphone
- Kulkas yang bisa memberitahu kalau makanan habis
- Dalam kasus kita: sensor tanah yang memberitahu kondisi nutrisi tanaman

---

## 🔧 Komponen-Komponen Sistem

### 1. **ESP32 - "Otak" Sensor**
```
ESP32 = Mikrokontroler (komputer mini)
```

**Apa itu ESP32?**
- Chip komputer kecil seukuran kotak korek api
- Bisa terhubung ke WiFi dan Bluetooth
- Bisa membaca data dari sensor
- Bisa mengirim data ke internet

**Fungsi dalam sistem kita:**
- Membaca data NPK (Nitrogen, Phosphorus, Kalium) dari tanah
- Membaca suhu dan kelembaban tanah
- Membaca tingkat pH tanah
- Mengirim semua data ke server melalui internet

**Contoh penggunaan:**
```javascript
// ESP32 membaca sensor setiap 30 detik
Setiap 30 detik:
  - Baca kadar Nitrogen: 25.5 mg/kg
  - Baca kadar Phosphorus: 15.2 mg/kg  
  - Baca kadar Kalium: 160.8 mg/kg
  - Baca suhu tanah: 28.5°C
  - Baca kelembaban: 75.2%
  - Baca pH: 6.8
  - Kirim semua data ke server
```

---

### 2. **MQTT - "Sistem Pos Surat" Digital**
```
MQTT = Message Queuing Telemetry Transport
```

**Apa itu MQTT?**
MQTT adalah protokol komunikasi yang sangat efisien untuk perangkat IoT. Bayangkan seperti sistem pos surat yang sangat cepat dan hemat energi.

**Mengapa pakai MQTT?**
- ✅ **Hemat energi** - Cocok untuk perangkat berbaterai
- ✅ **Cepat** - Data sampai dalam hitungan detik
- ✅ **Reliable** - Jika koneksi putus, data tidak hilang
- ✅ **Ringan** - Tidak memakan banyak bandwidth internet

**Cara kerja MQTT:**
```
1. PUBLISHER (ESP32) → Mengirim data sensor
2. MQTT BROKER → Seperti kantor pos yang menerima dan menyalurkan pesan
3. SUBSCRIBER (Server kita) → Menerima data sensor
```

**Topic dalam MQTT:**
```
Topic = Alamat tujuan pesan

Contoh topic yang kita gunakan:
- sugar_vestrack/sensor/ESP32_001/data     → Data sensor dari ESP32_001
- sugar_vestrack/device/ESP32_001/status   → Status perangkat ESP32_001
- sugar_vestrack/device/ESP32_001/command  → Perintah untuk ESP32_001
```

---

### 3. **NPK Sensor - "Mata" Pemantau Tanah**
```
NPK = Nitrogen + Phosphorus + Kalium (3 nutrisi utama tanaman)
```

**Apa itu Sensor NPK?**
Sensor khusus yang bisa mengukur kadar nutrisi dalam tanah, seperti dokter yang memeriksa kesehatan tanah.

**Yang diukur:**
- **Nitrogen (N)** - Untuk pertumbuhan daun
- **Phosphorus (P)** - Untuk pertumbuhan akar  
- **Kalium (K)** - Untuk ketahanan tanaman

**Nilai optimal untuk tanaman tebu:**
```
Nitrogen   : 20-30 mg/kg  (jika < 20 = kurang, jika > 40 = berlebih)
Phosphorus : 15-20 mg/kg  (jika < 15 = kurang, jika > 30 = berlebih)  
Kalium     : 150-180 mg/kg (jika < 150 = kurang, jika > 220 = berlebih)
pH         : 6.0-8.0      (jika < 6.0 atau > 8.0 = tidak optimal)
```

---

## 🏗️ Arsitektur Sistem (Cara Kerja Keseluruhan)

### Alur Data Sederhana:
```
1. Sensor NPK di tanah → Baca kondisi tanah
2. ESP32 → Kumpulkan data dari sensor
3. ESP32 → Kirim data via WiFi menggunakan MQTT
4. MQTT Broker → Terima dan salurkan data
5. Server Go → Terima data dan simpan ke database
6. Web Dashboard → Tampilkan data untuk petani
7. Sistem Alert → Kirim notifikasi jika ada masalah
```

### Diagram Sederhana:
```
[Tanah] → [Sensor NPK] → [ESP32] → [WiFi] → [MQTT Broker] → [Server] → [Database]
                                                              ↓
                                                         [Web Dashboard]
                                                         [Mobile App]
                                                         [Alert System]
```

---

## 💾 Database & Backend (Server)

### Apa yang Disimpan?
**Tabel sensor_data:**
```sql
Setiap data sensor berisi:
- ID perangkat (ESP32_001, ESP32_002, dst)
- Nama kebun (Farm_A, Farm_B, dst)
- Kadar Nitrogen, Phosphorus, Kalium
- Suhu dan kelembaban tanah
- Tingkat pH tanah
- Koordinat GPS lokasi sensor
- Waktu pengukuran
```

**Tabel device_status:**
```sql
Status setiap perangkat:
- Apakah perangkat online/offline
- Level baterai perangkat
- Kekuatan sinyal
- Versi firmware
- Terakhir terlihat kapan
```

**Tabel sensor_alerts:**
```sql
Alert yang dihasilkan sistem:
- Jenis alert (nitrogen_low, ph_abnormal, dst)
- Tingkat keparahan (low, medium, high, critical)
- Nilai sensor yang memicu alert
- Nilai batas yang dilanggar
- Status sudah ditangani atau belum
```

---

## 🚨 Sistem Alert Pintar

### Cara Kerja Alert:
```javascript
Sistem otomatis cek setiap data sensor masuk:

IF Nitrogen < 20 mg/kg:
  ➡️ Buat alert "Nitrogen rendah" - Level: Medium
  
IF pH < 5.0 atau pH > 9.0:
  ➡️ Buat alert "pH abnormal" - Level: Critical
  
IF Suhu > 40°C:
  ➡️ Buat alert "Suhu terlalu tinggi" - Level: High
```

### Level Keparahan:
- 🟢 **Low** - Perlu perhatian tapi tidak urgent
- 🟡 **Medium** - Perlu tindakan dalam beberapa hari
- 🟠 **High** - Perlu tindakan segera
- 🔴 **Critical** - Perlu tindakan darurat

---

## 🌐 REST API (Interface Komunikasi)

### Apa itu REST API?
REST API adalah cara untuk aplikasi web/mobile berkomunikasi dengan server. Seperti menu restoran yang menjelaskan makanan apa saja yang tersedia.

### Endpoint yang Tersedia:

**1. Ambil Data Sensor:**
```
GET /api/sensor/data
→ Ambil semua data sensor dengan filter tertentu

Contoh response:
{
  "status": "success",
  "data": [
    {
      "device_id": "ESP32_001",
      "farm_name": "Farm_A", 
      "nitrogen": 25.5,
      "phosphorus": 15.2,
      "potassium": 160.8,
      "temperature": 28.5,
      "humidity": 75.2,
      "ph": 6.8,
      "timestamp": "2025-06-27T10:30:00Z"
    }
  ]
}
```

**2. Cek Status Perangkat:**
```
GET /api/sensor/devices/status
→ Lihat perangkat mana yang online/offline

Contoh response:
{
  "status": "success",
  "data": [
    {
      "device_id": "ESP32_001",
      "farm_name": "Farm_A",
      "is_online": true,
      "battery_level": 85.5,
      "last_seen": "2025-06-27T15:45:00Z"
    }
  ]
}
```

**3. Kirim Perintah ke Perangkat:**
```
POST /api/sensor/devices/ESP32_001/command
Body: {
  "command": "restart",
  "payload": {
    "reason": "maintenance"
  }
}
→ Restart perangkat ESP32_001 untuk maintenance
```

**4. Lihat Alert:**
```
GET /api/sensor/alerts
→ Lihat semua alert yang belum ditangani

Contoh response:
{
  "status": "success", 
  "data": [
    {
      "device_id": "ESP32_001",
      "alert_type": "nitrogen_low",
      "message": "Nitrogen level too low: 18.5 mg/kg",
      "severity": "medium",
      "is_resolved": false
    }
  ]
}
```

---

## 🛠️ Teknologi yang Digunakan

### Backend (Server):
- **Go (Golang)** - Bahasa pemrograman untuk server
- **Gin Framework** - Framework web untuk membuat REST API
- **GORM** - Library untuk mengakses database
- **Paho MQTT** - Library untuk komunikasi MQTT
- **MySQL** - Database untuk menyimpan data

### Frontend (Tampilan):
- **Swagger UI** - Dokumentasi interaktif API
- **HTML/CSS/JavaScript** - Untuk web dashboard (bisa dikembangkan)

### Hardware:
- **ESP32** - Mikrokontroler
- **Sensor NPK** - Sensor nutrisi tanah
- **Sensor Suhu/Kelembaban** - Sensor kondisi lingkungan

---

## 📊 Manfaat Sistem untuk Petani

### 1. **Monitoring Real-time**
```
Petani bisa lihat kondisi tanah setiap saat dari smartphone:
- "Kebun A butuh pupuk nitrogen"
- "Kebun B pH tanah terlalu asam" 
- "Kebun C kondisi optimal"
```

### 2. **Alert Otomatis**
```
Sistem otomatis memberitahu jika ada masalah:
- SMS: "ALERT - Nitrogen rendah di Kebun A, segera beri pupuk"
- WhatsApp: "pH tanah Kebun B terlalu tinggi, perlu treatment"
```

### 3. **Data Historis**
```
Petani bisa lihat tren:
- "Bulan lalu nitrogen selalu turun setelah hujan"
- "Kebun ini selalu pH tinggi di musim kemarau"
- "Pola terbaik untuk panen adalah..."
```

### 4. **Efisiensi Pupuk**
```
Tidak perlu tebak-tebakan:
- Pupuk nitrogen hanya diberikan saat benar-benar kurang
- Tidak membuang pupuk berlebihan
- Hemat biaya operasional
```

### 5. **Remote Monitoring**
```
Petani tidak perlu ke kebun setiap hari:
- Pantau dari rumah via smartphone
- Dapat notifikasi otomatis
- Fokus ke kebun yang bermasalah saja
```

---

## 🚀 Cara Menjalankan Sistem

### Persiapan:
1. **Setup Hardware:**
   - Pasang ESP32 + Sensor NPK di setiap kebun
   - Pastikan ada koneksi WiFi
   - Atur power supply (baterai/solar panel)

2. **Setup Software:**
   - Install aplikasi server di komputer/cloud
   - Setup database MySQL
   - Konfigurasi MQTT broker

3. **Testing:**
   - Jalankan ESP32 simulator untuk test
   - Cek data masuk ke database
   - Test alert system
   - Test web dashboard

### Menjalankan:
```bash
# 1. Jalankan server utama
./bin/main.exe

# 2. Jalankan simulator ESP32 (untuk testing)
./tools/esp32_simulator.exe

# 3. Buka browser ke:
http://localhost:8080/swagger/index.html
```

---

## 🔍 Troubleshooting Umum

### Masalah Koneksi:
```
Problem: ESP32 tidak terkoneksi
Solution: 
- Cek koneksi WiFi
- Cek konfigurasi MQTT broker
- Restart ESP32
```

### Masalah Data:
```
Problem: Data sensor tidak masuk
Solution:
- Cek sensor fisik
- Cek kabel koneksi
- Cek log error di server
```

### Masalah Alert:
```
Problem: Alert tidak muncul
Solution:
- Cek threshold setting
- Cek service alert berjalan
- Cek notifikasi setting
```

---

## 📈 Pengembangan Selanjutnya

### Fitur yang Bisa Ditambah:

1. **Mobile App Native:**
   - Android/iOS app
   - Push notification
   - Offline capability

2. **Advanced Analytics:**
   - Machine learning prediction
   - Weather integration
   - Yield forecasting

3. **Additional Sensors:**
   - Soil moisture sensor
   - Light intensity sensor
   - Pest detection camera

4. **Automation:**
   - Auto irrigation system
   - Auto fertilizer dispenser
   - Drone integration

5. **Business Intelligence:**
   - Cost analysis dashboard
   - ROI calculation
   - Market price integration

---

## 💡 Kesimpulan

Sistem IoT monitoring tanaman tebu ini memberikan solusi modern untuk pertanian tradisional. Dengan teknologi MQTT dan ESP32, petani dapat:

- ✅ **Memantau kondisi tanah secara real-time**
- ✅ **Mendapat alert otomatis saat ada masalah**
- ✅ **Menghemat biaya pupuk dan perawatan**
- ✅ **Meningkatkan produktivitas tanaman**
- ✅ **Membuat keputusan berdasarkan data akurat**

Sistem ini menggabungkan hardware (ESP32, sensor), software (Go backend), protokol komunikasi (MQTT), dan database untuk menciptakan solusi IoT yang lengkap dan mudah digunakan.

---

**Dibuat dengan ❤️ untuk kemajuan pertanian Indonesia**
*Dokumen ini ditulis dengan bahasa sederhana agar mudah dipahami oleh siapa saja yang ingin belajar tentang teknologi IoT dalam pertanian.*
