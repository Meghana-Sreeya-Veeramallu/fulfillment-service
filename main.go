package main

import (
	"Fulfillment/Controller"
	"Fulfillment/Model"
	"Fulfillment/Repository"
	"Fulfillment/Service"
	pb "Fulfillment/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
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

	// Create the repository and service
	repo := Repository.NewDeliveryAgentRepository(db)
	deliveryAgentService := Service.NewDeliveryAgentService(repo)

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Create the controller and pass the service
	deliveryAgentServer := Controller.NewDeliveryAgentServer(deliveryAgentService)
	pb.RegisterDeliveryAgentServiceServer(s, deliveryAgentServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
