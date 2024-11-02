package Controller

import (
	"Fulfillment/Model"
	pb "Fulfillment/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *gorm.DB
var lis *bufconn.Listener
var server *grpc.Server
var client pb.DeliveryAgentServiceClient

const bufSize = 1024 * 1024 // 1MB

// Setup function to initialize the test database and server
func setup() {
	db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	err := db.AutoMigrate(&Model.DeliveryAgent{})
	if err != nil {
		panic(err)
	}

	lis = bufconn.Listen(bufSize)
	server = grpc.NewServer()
	pb.RegisterDeliveryAgentServiceServer(server, NewDeliveryAgentServer(db))

	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(
		func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	client = pb.NewDeliveryAgentServiceClient(conn)
}

// Teardown function to stop the server
func teardown() {
	server.Stop()
}

// Test AddDeliveryAgent
func TestAddDeliveryAgentSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	req := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: "Hyderabad"}
	resp, err := client.AddDeliveryAgent(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Delivery agent added successfully", resp.Message)
}

// Test AddDeliveryAgent with invalid name
func TestAddDeliveryAgentInvalidName(t *testing.T) {
	setup()
	defer teardown()

	req := &pb.AddDeliveryAgentRequest{Name: "", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "name cannot be empty", status.Convert(err).Message())
}

// Test AddDeliveryAgent with invalid city
func TestAddDeliveryAgentInvalidCity(t *testing.T) {
	setup()
	defer teardown()

	req := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: ""}
	_, err := client.AddDeliveryAgent(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "city cannot be empty", status.Convert(err).Message())
}

// Test AssignAgentToOrder
func TestAssignAgentToOrderSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	addReq := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), addReq)
	assert.NoError(t, err)

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 1, OrderId: 123}
	resp, err := client.AssignAgentToOrder(context.Background(), assignReq)
	assert.NoError(t, err)
	assert.Equal(t, "Delivery agent assigned to order successfully", resp.Message)
}

// Test AssignAgentToOrder when agent not found
func TestAssignAgentToOrderWhenDeliveryAgentNotFound(t *testing.T) {
	setup()
	defer teardown()

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 999, OrderId: 123}
	_, err := client.AssignAgentToOrder(context.Background(), assignReq)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "delivery agent not found", status.Convert(err).Message())
}

// Test AssignAgentToOrder when already assigned
func TestAssignAgentToOrderWhenAlreadyAssigned(t *testing.T) {
	setup()
	defer teardown()

	addReq := &pb.AddDeliveryAgentRequest{Name: "Harish", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), addReq)
	assert.NoError(t, err)

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 1, OrderId: 123}
	_, err = client.AssignAgentToOrder(context.Background(), assignReq)
	assert.NoError(t, err)

	_, err = client.AssignAgentToOrder(context.Background(), assignReq)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "delivery agent is not available", status.Convert(err).Message())
}
