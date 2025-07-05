package db

import (
	"customer-voucher-service/models/brand_model"
	"customer-voucher-service/models/customer_model"
	"customer-voucher-service/models/transaction_model"
	"customer-voucher-service/models/voucher_model"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	fmt.Println("Using DB connection:")
	fmt.Println("Host:", os.Getenv("DB_HOST"))
	fmt.Println("User:", os.Getenv("DB_USER"))
	fmt.Println("DB:", os.Getenv("DB_NAME"))
	fmt.Println("Port:", os.Getenv("DB_PORT"))
	fmt.Println("SSL Mode:", os.Getenv("DB_SSLMODE"))

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	err = DB.AutoMigrate(
		&brand_model.Brand{},
		&voucher_model.Voucher{},
		&customer_model.Customer{},
		&transaction_model.Transaction{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate:", err)
	}
}
