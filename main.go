package main

import (
	"Fulfillment/Controller"
	"Fulfillment/Model"
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

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	deliveryAgentServer := Controller.NewDeliveryAgentServer(db)
	pb.RegisterDeliveryAgentServiceServer(s, deliveryAgentServer)

	log.Printf("Server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
