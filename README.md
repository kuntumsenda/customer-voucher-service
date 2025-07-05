# customer-voucher-service

## Setup Database

Sebelum menjalankan aplikasi, pastikan database `voucher_db` sudah dibuat di PostgreSQL. Aplikasi **tidak** akan membuat database secara otomatis, hanya tabel-tabel di dalam database yang sudah ada.

### Cara Membuat Database

1. Masuk ke PostgreSQL menggunakan terminal atau tools seperti DBeaver.
2. Jalankan perintah berikut di terminal PostgreSQL:

   ```sql
   CREATE DATABASE voucher_db;
   ```

3. Pastikan user, host, dan port sesuai dengan konfigurasi di file `.env` atau environment variable aplikasi.

Setelah database `voucher_db` tersedia, jalankan aplikasi dengan:

```bash
go run main.go
```

Aplikasi akan otomatis membuat tabel-tabel yang diperlukan di dalam database `voucher_db`.

