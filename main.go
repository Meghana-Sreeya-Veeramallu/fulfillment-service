package main

import (
	"Fulfillment/DeliveryAgent"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=postgres password=admin dbname=fulfillment_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	if err := db.AutoMigrate(&DeliveryAgent.DeliveryAgent{}); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}
