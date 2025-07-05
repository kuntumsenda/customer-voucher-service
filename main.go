package main

import (
	"customer-voucher-service/db"
	"customer-voucher-service/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	r := gin.Default()

	routes.ApiRoutes(r)

	r.Run(":8080")
}
