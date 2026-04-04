# SiCemas - Sistem Informasi Catatan Elektronik Masjid Area Ciemas

Proyek ini merupakan hasil penelitian Tugas Akhir (Skripsi) saya untuk program studi Informatika di UNIKOM. Fokus utama dari pengembangan SiCemas adalah membangun platform manajemen data situs keagamaan (seperti Masjid, Pesantren, Mushola, dan Majelis Taklim) yang tidak hanya fungsional secara fitur, tetapi juga memiliki standar keamanan yang ketat sesuai dengan kebutuhan instansi pemerintah.

Dalam proyek ini, saya mengimplementasikan arsitektur keamanan **Zero Trust** dan **Column-Level Encryption** untuk memastikan integritas data dan perlindungan privasi pengguna di lingkungan Kantor Urusan Agama (KUA) Kecamatan Ciemas.

## Apa yang Spesial di Proyek Ini?

Berbeda dengan aplikasi manajemen data biasa, SiCemas dirancang dengan pendekatan _security-first_:

1. **Arsitektur Zero Trust**: Menggunakan prinsip _"Never Trust, Always Verify"_. Sistem menerapkan pengecekan _Trust Score_ secara kontinyu melalui middleware yang memvalidasi IP Address (Geo-IP via MaxMind), sidik jari perangkat (_device fingerprinting_), dan konsistensi sesi. Jika terdeteksi anomali pada konteks akses, sesi akan otomatis dicabut (invalidated).
2. **Column-Level Encryption**: Melindungi data sensitif pengguna dengan enkripsi AES pada level kolom database PostgreSQL. Hal ini memitigasi risiko kebocoran data (data breach) pada layer penyimpanan.
3. **Single Binary Deployment**: Memanfaatkan fitur `go:embed` untuk menyatukan seluruh aset frontend SvelteKit yang sudah dikompilasi ke dalam satu file binary Go. Aplikasi dapat dijalankan di server produksi tanpa ketergantungan pada Node.js atau folder aset eksternal.
4. **Cloud Storage Integration**: Manajemen foto situs keagamaan menggunakan integrasi Cloudinary API untuk penyimpanan gambar yang teroptimasi dan aman.

## Teknologi yang Digunakan

### Backend

- **Go (Golang)**: Bahasa pemrograman utama untuk skalabilitas dan performa tinggi.
- **Go Fiber**: Framework web minimalis untuk menangani routing API.
- **SQLx & PostgreSQL**: Manajemen database relasional dengan performa query yang optimal.
- **Redis**: Digunakan sebagai penyimpanan sementara untuk session management, rate limiting, dan caching OTP.
- **Casbin**: Implementasi Role-Based Access Control (RBAC) untuk manajemen hak akses petugas.

### Frontend

- **Svelte 5 (Runes)**: Mengadopsi fitur terbaru Svelte untuk reaktivitas state yang lebih efisien.
- **SvelteKit**: Meta-framework untuk routing dan optimalisasi Single Page Application (SPA).
- **TypeScript**: Memastikan keamanan tipe data (type-safety) di seluruh komponen frontend.
- **Valibot**: Validasi skema form yang ringan dan terintegrasi dengan TypeScript.
- **TailwindCSS & DaisyUI**: Framework CSS untuk antarmuka yang modern dan responsif.

## Struktur Direktori

- `cmd/api/`: Inisialisasi server, konfigurasi environment, dan dependency injection.
- `internal/app/`: Implementasi logic bisnis yang terbagi menjadi Handlers, Services, dan Repositories.
- `internal/middlewares/`: Layer keamanan (Auth, Zero Trust Validator, Rate Limiter).
- `internal/utils/`: Helper fungsional termasuk enkripsi AES, JWT, dan pengolahan data Excel.
- `ui/`: Source code frontend berbasis SvelteKit.
- `migrations/`: Skema database PostgreSQL yang mendukung mekanisme soft-delete.

## Cara Menjalankan Aplikasi

### Menjalankan Secara Lokal

1. Pastikan sudah menginstal Go (v1.22+), Bun/Node.js, PostgreSQL, dan Redis.
2. Clone repositori ini:
   ```bash
   git clone https://github.com/Auxesia23/sicemas.git
   ```
3. Konfigurasi environment:
   ```bash
   cp .env.example .env
   ```
   Sesuaikan detail DB, Redis, Cloudinary, dan SMTP.
4. Persiapkan Frontend:
   ```bash
   cd ui && bun install && bun run build && cd ..
   ```
5. Jalankan migrasi database dan aplikasi:
   ```bash
   make migrate-up
   make run
   ```

### Menggunakan Docker Compose

Untuk cara yang lebih praktis dan terisolasi, Anda dapat menggunakan Docker Compose. Metode ini akan menjalankan aplikasi, database, dan redis secara otomatis dalam satu jaringan.

1. Persiapkan environment
   Pastikan file .env sudah sesuai (gunakan nama service db untuk host database dan redis untuk host redis).
2. Build dan jalankan service
   Gunakan perintah berikut di root direktori project:
   ```bash
   docker compose up -d --build
   ```
3. Verifikasi service
   Cek apakah semua kontainer sudah berjalan:
   ```bash
   docker compose ps
   ```

---

**Fajar Gustiana (Auxesia23)** Teknik Informatika - UNIKOM Bandung
