# Customer Voucher Service

A Go-based microservice for managing customer vouchers, brands, and transactions.

## Prerequisites

- Go 1.19 or higher
- PostgreSQL 12 or higher
- Protocol Buffers compiler (protoc)

## Setup

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd customer-voucher-service

# Install Go dependencies
go mod tidy
```

### 2. Setup Database

Sebelum menjalankan aplikasi, pastikan database `voucher_db` sudah dibuat di PostgreSQL. Aplikasi **tidak** akan membuat database secara otomatis, hanya tabel-tabel di dalam database yang sudah ada.

#### Cara Membuat Database

1. Masuk ke PostgreSQL menggunakan terminal atau tools seperti DBeaver (OWNER used this).
2. Jalankan perintah berikut di terminal PostgreSQL:

   ```sql
   CREATE DATABASE voucher_db;
   ```

3. Pastikan user, host, dan port sesuai dengan konfigurasi di file `.env` atau environment variable aplikasi.

### 3. Environment Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=voucher_db
DB_SSLMODE=disable
```

### 4. Generate Protocol Buffers

```bash
# Generate Go code from .proto files
make proto
```

### 5. Run the Application

```bash
go run main.go
```

Aplikasi akan otomatis membuat tabel-tabel yang diperlukan di dalam database `voucher_db`.

## API Endpoints

The service provides the following main endpoints:

- **Customer Management**: Create, list, and manage customers
- **Brand Management**: Create and list brands
- **Voucher Management**: Create, list, and manage vouchers
- **Transaction Management**: Redeem points, list transactions, and view transaction details

## Testing

Run unit tests:

```bash
# Run all tests
go test ./...

# Run specific service tests
go test ./services/customer_service/ -v
go test ./services/brand_service/ -v
go test ./services/voucher_service/ -v
go test ./services/transaction_service/ -v

# Run model tests
go test ./models/voucher_model/ -v
```

## Development

### Adding New Dependencies

```bash
go get <package-name>
go mod tidy
```

### Regenerating Protocol Buffers

After modifying `.proto` files:

```bash
make proto
```

### Database Migrations

Aplikasi ini menggunakan auto-migration GORM. Saat menambahkan model baru di file entity, model tersebut akan secara otomatis dibuat dalam database saat dijalankan.

