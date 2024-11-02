package Controller

import (
	"Fulfillment/Service"
	pb "Fulfillment/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

// DeliveryAgentServer represents the gRPC server for delivery agents.
type DeliveryAgentServer struct {
	pb.UnimplementedDeliveryAgentServiceServer
	DB *gorm.DB
}

// NewDeliveryAgentServer creates a new DeliveryAgentServer.
func NewDeliveryAgentServer(db *gorm.DB) *DeliveryAgentServer {
	return &DeliveryAgentServer{DB: db}
}

// AddDeliveryAgent handles the gRPC request for adding a delivery agent.
func (s *DeliveryAgentServer) AddDeliveryAgent(ctx context.Context, req *pb.AddDeliveryAgentRequest) (*pb.AddDeliveryAgentResponse, error) {
	_, err := Service.AddDeliveryAgent(s.DB, req.Name, req.City)
	if err != nil {
		if err.Error() == "name cannot be empty" || err.Error() == "city cannot be empty" {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, grpc.Errorf(codes.Internal, "Failed to add delivery agent: "+err.Error())
	}

	return &pb.AddDeliveryAgentResponse{
		Message: "Delivery agent added successfully",
	}, nil
}

// AssignAgentToOrder handles the gRPC request for assigning an agent to an order.
func (s *DeliveryAgentServer) AssignAgentToOrder(ctx context.Context, req *pb.AssignAgentToOrderRequest) (*pb.AssignAgentToOrderResponse, error) {
	agentID := req.AgentId
	orderID := req.OrderId

	err := Service.AssignAgentToOrder(s.DB, uint(agentID), int(orderID))
	if err != nil {
		if err.Error() == "delivery agent not found" || err.Error() == "delivery agent is not available" {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, grpc.Errorf(codes.Internal, "An error occurred: "+err.Error())
	}

	return &pb.AssignAgentToOrderResponse{
		Message: "Delivery agent assigned to order successfully",
	}, nil
}
