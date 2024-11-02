package main

import (
	"Fulfillment/Controller"
	"Fulfillment/Model"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	dsn := "host=localhost user=postgres password=admin dbname=fulfillment_service port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	if err := db.AutoMigrate(&Model.DeliveryAgent{}); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	router := mux.NewRouter()
	deliveryAgentServer := Controller.NewDeliveryAgentServer(db)
	deliveryAgentServer.RegisterRoutes(router)

	log.Println("Starting HTTP server on port 8082...")
	if err := http.ListenAndServe(":8082", router); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
