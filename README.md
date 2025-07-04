
# Starter Kit Backend Golang (swagger)

## Deskripsi
Starter Kit Backend Golang adalah kit awal untuk memulai pengembangan aplikasi backend menggunakan bahasa Go. Proyek ini menyediakan struktur modular dengan komponen-komponen yang siap digunakan seperti:

- **Controllers**: Mengelola logika bisnis seperti autentikasi, kategori, produk, dan izin.
- **Helpers**: Fungsi bantuan untuk pengelolaan lingkungan, hashing, dan respons.
- **Middleware**: Pengelolaan autentikasi dan logging.
- **Models**: Definisi skema data (pengguna, produk, kategori).
- **Services**: Logika bisnis aplikasi.
- **Routes**: Pengaturan rute aplikasi.
- **Database**: Mengelola struktur dan data awal database, termasuk migrasi, seeder, serta CLI untuk pengelolaan skema data.

Proyek ini juga terintegrasi dengan Swagger untuk dokumentasi API. Kamu bisa mengakses dokumentasi Swagger melalui Swagger UI setelah menjalankan proyek.

## Fitur Utama
- Autentikasi menggunakan JWT
- Manajemen pengguna dan otoritas
- Hot-reload dengan Air
- Dokumentasi API otomatis dengan Swagger
- Arsitektur modular dengan layanan terpisah
- Sistem migrasi database: migrate, rollback, fresh, make:migration
- Sistem seeder data awal: db:seed, rollback:seeder, make:seeder

## Konfigurasi Environment
Proyek ini menggunakan file `.env` untuk mengatur konfigurasi. Berikut adalah contoh file `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
JWT_SECRET=your_jwt_secret_key
```

## Cara Menjalankan Aplikasi

1. **Clone repository**:
    ```bash
    git clone https://github.com/RahmatRafiq/golang_starter_kit_2025.git
    ```
2. **Masuk ke direktori project**:
    ```bash
    cd golang_starter_kit_2025
    ```
3. **Instal dependencies**:
    ```bash
    go install github.com/air-verse/air@lates
    go install github.com/swaggo/swag/cmd/swag@latest
    ```
4. **Generate dokumentasi Swagger**:
    ```bash
    swag init
    ```
5. **Jalankan aplikasi dengan Air**:
    ```bash
    air
    ```
6. **Akses halaman API dengan URL**:
    ```bash
    http://localhost:8080/swagger/index.html
    ```

## Perintah CLI untuk Migrasi

### 1. Membuat File Migrasi Baru
```bash
go run main.go make:migration <prefix_nama_migrasi>
```
Membuat satu file kosong dengan format:
- `YYYYMMDDHHMMSS_<prefix_nama_migrasi>.sql`

📌 **Rekomendasi**:
- Gunakan prefix seperti `create_` atau `alter_` untuk mempermudah identifikasi jenis migrasi.
- Contoh:
    - `create_users_table`
    - `alter_products_table`

### 2. Menjalankan Satu File Migrasi
```bash
go run main.go migrate --file <nama_file_migration>
```
Menjalankan bagian `UP` dari `<nama_file_migration>.sql`.

### 3. Menjalankan Semua Migrasi yang Tertunda
```bash
go run main.go migrate:all
```
- Membuat batch baru.
- Menjalankan bagian `UP` dari semua file `.sql` yang belum tercatat di tabel `migrations`.
- Mencatat setiap file ke batch tersebut.

### 4. Rollback Satu File Migrasi
```bash
go run main.go rollback --file=<nama_file_migration>
```
Menjalankan bagian `DOWN` dari `<nama_file_migration>.sql` (tanpa mengubah batch).

### 5. Rollback Semua Batch
```bash
go run main.go rollback:all
```
- Loop dari batch tertinggi → 1.
- Menjalankan bagian `DOWN` dari semua file `.sql` per batch.
- Menghapus seluruh catatan di tabel `migrations`.

### 6. Rollback Batch Tertentu
```bash
go run main.go rollback:batch --batch=<nomor_batch>
```
Meng-rollback hanya migrasi di batch `<nomor_batch>`, lalu menghapus catatannya.

### 7. Rollback Batch Terakhir (Default)
```bash
go run main.go rollback:batch
```
Jika flag `--batch` tidak diset, akan otomatis meng-rollback batch terakhir.

---

📌 **Catatan**:
- Tabel `migrations` akan otomatis dibuat saat pertama kali menjalankan `migrate:all` atau `rollback:batch`.
- Pastikan setiap file `.sql` memiliki bagian `UP` dan `DOWN` yang jelas sebelum menjalankan migrate/rollback.
- Contoh format file migrasi:
    ```sql
    -- UP
    CREATE TABLE users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    );

    -- DOWN
    DROP TABLE users;
    ```


## Perintah CLI untuk Seeder

### 1. Membuat File Seeder Baru
```bash
go run main.go make:seeder --name=<nama_seeder>
```
Membuat file seeder baru di direktori `app/database/seeds/` dengan format nama file:
- `YYYYMMDDHHMMSS_<nama_seeder>.go`

### 2. Menjalankan Semua Seeder
```bash
go run main.go db:seed
```
Menjalankan semua seeder yang ada di direktori `app/database/seeds/`.

### 3. Rollback Batch Seeder Terakhir (Default)
```bash
go run main.go rollback:seeder
```
Menghapus data yang dimasukkan oleh batch seeder terakhir.

### 4. Rollback Batch Seeder Tertentu
```bash
go run main.go rollback:seeder --batch=<nomor_batch>
```
Menghapus data yang dimasukkan oleh batch seeder dengan nomor `<nomor_batch>`.

📌 **Catatan**:
- Seeder file yang dibuat akan memiliki template dasar untuk mempermudah implementasi.
- Pastikan untuk menyesuaikan isi file seeder dengan kebutuhan data aplikasi Anda.
- Gunakan perintah rollback untuk menghapus data yang tidak diperlukan atau untuk pengujian ulang.

Semoga bermanfaat! 😊

## Kontribusi

Aplikasi ini dikembangkan oleh [Dzyfhuba](https://github.com/Dzyfhuba) dan [RahmatRafiq](https://github.com/RahmatRafiq). Jangan ragu untuk mengajukan pertanyaan atau memberikan saran. Kami sangat terbuka terhadap kontribusi dari siapa saja yang ingin terlibat.

<div align="center">

# Selamat mencoba, dan semoga proyek ini membantu kamu dalam pengembangan aplikasi backend! 😊

</div>


Terima kasih atas dukunganmu! 🙏


## Struktur Proyek
Berikut adalah struktur direktori proyek ini beserta deskripsi singkatnya:
```markdown
golang_starter_kit_2025/
├── app/
│   ├── **controllers/** - Mengelola logika bisnis aplikasi
│   │   ├── `authController.go` - Logika autentikasi
│   │   ├── `productController.go` - Logika produk
│   │   └── `userController.go` - Logika pengguna
│   ├── **helpers/** - Fungsi bantuan untuk berbagai kebutuhan
│   │   ├── `environmentHelper.go` - Pengelolaan variabel lingkungan
│   │   ├── `responseHelper.go` - Format respons API
│   │   └── `hashingHelper.go` - Fungsi hashing
│   ├── **middleware/** - Middleware untuk aplikasi
│   │   ├── `authMiddleware.go` - Middleware autentikasi
│   │   ├── `loggingMiddleware.go` - Middleware logging
│   │   └── `errorHandlerMiddleware.go` - Middleware penanganan error
│   ├── **models/** - Definisi skema data
│   │   ├── `userModel.go` - Model pengguna
│   │   ├── `productModel.go` - Model produk
│   │   └── `categoryModel.go` - Model kategori
│   ├── **services/** - Logika bisnis aplikasi
│   │   ├── `authService.go` - Layanan autentikasi
│   │   ├── `productService.go` - Layanan produk
│   │   └── `userService.go` - Layanan pengguna
│   ├── **routes/** - Pengaturan rute aplikasi
│   │   ├── `authRoutes.go` - Rute autentikasi
│   │   ├── `productRoutes.go` - Rute produk
│   │   └── `userRoutes.go` - Rute pengguna
│   ├── **database/** - Pengelolaan database
│   │   ├── **migrations/** - File migrasi database
│   │   │   ├── `20230424010101_create_users_table.up.sql` - Membuat tabel pengguna
│   │   │   ├── `20230424010101_create_users_table.down.sql` - Menghapus tabel pengguna
│   │   │   └── `20230424010202_create_products_table.up.sql` - Membuat tabel produk
│   │   ├── **seeds/** - File seeder data awal
│   │   │   ├── `20230424010101_users_seeder.go` - Seeder pengguna
│   │   │   └── `20230424010202_products_seeder.go` - Seeder produk
│   │   └── `database.go` - Koneksi dan konfigurasi database
├── **config/** - Konfigurasi aplikasi
│   ├── `config.go` - Konfigurasi utama
│   └── `env.go` - Pengelolaan variabel lingkungan
├── **docs/** - Dokumentasi API
│   └── **swagger/** - File Swagger
│       └── `swagger.yaml` - Dokumentasi API Swagger
├── **bootstrap/** - Inisialisasi aplikasi
│   ├── `main.go` - Entry point aplikasi
│   └── `init.go` - Inisialisasi komponen aplikasi
├── `.env` - File konfigurasi lingkungan
├── `go.mod` - File modul Go
├── `go.sum` - File checksum dependensi
└── `README.md` - Dokumentasi proyek
```

Struktur ini dirancang untuk mempermudah pengelolaan kode, meningkatkan keterbacaan, dan mendukung pengembangan aplikasi secara modular.


## 💸 Dukung Proyek Ini

> 🇮🇩 [Donasi via Saweria (Indonesia)](https://saweria.co/RahmatRafiq)  