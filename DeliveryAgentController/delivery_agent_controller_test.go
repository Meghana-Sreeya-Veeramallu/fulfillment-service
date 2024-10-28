package DeliveryAgentController

import (
	"Fulfillment/DeliveryAgent"
	pb "Fulfillment/proto"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *gorm.DB
var router *mux.Router

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	err := db.AutoMigrate(&DeliveryAgent.DeliveryAgent{})
	if err != nil {
		return nil
	}
	return db
}

func setupTestServer(db *gorm.DB) *DeliveryAgentServer {
	return NewDeliveryAgentServer(db)
}

func setup(t *testing.T) {
	db = setupTestDB()
	if db == nil {
		t.Fatal("failed to create test database")
	}

	server := setupTestServer(db)
	router = mux.NewRouter()
	server.RegisterRoutes(router)
}

func teardown() {
	db, _ = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	db.Exec("DROP TABLE delivery_agents")
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestAddDeliveryAgentSuccessfully(t *testing.T) {
	setup(t)
	defer teardown()

	reqBody, _ := json.Marshal(&pb.AddDeliveryAgentRequest{
		Name: "Ketan",
		City: "Hyderabad",
	})
	req, _ := http.NewRequest("POST", "/delivery-agents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp pb.AddDeliveryAgentResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(t, "Delivery agent added successfully", resp.Message)
}

func TestAddDeliveryAgentInvalidName(t *testing.T) {
	setup(t)
	defer teardown()

	reqBody, _ := json.Marshal(&pb.AddDeliveryAgentRequest{
		Name: "",
		City: "Hyderabad",
	})
	req, _ := http.NewRequest("POST", "/delivery-agents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "name cannot be empty", strings.TrimSpace(rr.Body.String()))
}

func TestAddDeliveryAgentInvalidCity(t *testing.T) {
	setup(t)
	defer teardown()

	reqBody, _ := json.Marshal(&pb.AddDeliveryAgentRequest{
		Name: "Ketan",
		City: "",
	})
	req, _ := http.NewRequest("POST", "/delivery-agents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "city cannot be empty", strings.TrimSpace(rr.Body.String()))
}

func TestAssignAgentToOrderSuccessfully(t *testing.T) {
	setup(t)
	defer teardown()

	reqBody, _ := json.Marshal(&pb.AddDeliveryAgentRequest{
		Name: "Ketan",
		City: "Hyderabad",
	})
	req, _ := http.NewRequest("POST", "/delivery-agents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var addResp pb.AddDeliveryAgentResponse
	json.Unmarshal(rr.Body.Bytes(), &addResp)
	assert.Equal(t, "Delivery agent added successfully", addResp.Message)

	var addedAgent DeliveryAgent.DeliveryAgent
	if err := db.First(&addedAgent, uint32(1)).Error; err != nil {
		t.Fatalf("Failed to find added agent: %v", err)
	}
	assert.NotNil(t, addedAgent)

	assignReq, _ := http.NewRequest("POST", "/delivery-agents/1/orders/123", nil)
	assignReq.Header.Set("Content-Type", "application/json")

	assignRR := httptest.NewRecorder()
	router.ServeHTTP(assignRR, assignReq)

	assert.Equal(t, http.StatusOK, assignRR.Code)

	var assignResp pb.AssignAgentToOrderResponse
	json.Unmarshal(assignRR.Body.Bytes(), &assignResp)
	assert.Equal(t, "Delivery agent assigned to order successfully", assignResp.Message)

	var updatedAgent DeliveryAgent.DeliveryAgent
	db.First(&updatedAgent, uint32(1))

	orderID := 123
	assert.Equal(t, orderID, *updatedAgent.OrderID)
	assert.Equal(t, DeliveryAgent.UNAVAILABLE, updatedAgent.AvailabilityStatus)
}

func TestAssignAgentToOrderWhenDeliveryAgentNotFound(t *testing.T) {
	setup(t)
	defer teardown()

	assignReq, _ := http.NewRequest("POST", "/delivery-agents/999/orders/123", nil)
	assignReq.Header.Set("Content-Type", "application/json")

	assignRR := httptest.NewRecorder()
	router.ServeHTTP(assignRR, assignReq)

	assert.Equal(t, http.StatusBadRequest, assignRR.Code)
	assert.Equal(t, "delivery agent not found", strings.TrimSpace(assignRR.Body.String()))
}

func TestAssignAgentToOrderWhenAlreadyAssigned(t *testing.T) {
	setup(t)
	defer teardown()

	reqBody, _ := json.Marshal(&pb.AddDeliveryAgentRequest{
		Name: "Harish",
		City: "Hyderabad",
	})
	req, _ := http.NewRequest("POST", "/delivery-agents", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var addedAgent DeliveryAgent.DeliveryAgent
	if err := db.First(&addedAgent, uint32(1)).Error; err != nil {
		t.Fatalf("Failed to find added agent: %v", err)
	}
	assert.NotNil(t, addedAgent)

	assignReq, _ := http.NewRequest("POST", "/delivery-agents/1/orders/123", nil)
	assignReq.Header.Set("Content-Type", "application/json")

	assignRR := httptest.NewRecorder()
	router.ServeHTTP(assignRR, assignReq)

	assert.Equal(t, http.StatusOK, assignRR.Code)

	var updatedAgent DeliveryAgent.DeliveryAgent
	db.First(&updatedAgent, uint32(1))

	orderID := 123
	assert.Equal(t, orderID, *updatedAgent.OrderID)
	assert.Equal(t, DeliveryAgent.UNAVAILABLE, updatedAgent.AvailabilityStatus)

	assignReq2, _ := http.NewRequest("POST", "/delivery-agents/1/orders/456", nil)
	assignReq2.Header.Set("Content-Type", "application/json")

	assignRR2 := httptest.NewRecorder()
	router.ServeHTTP(assignRR2, assignReq2)

	assert.Equal(t, http.StatusBadRequest, assignRR2.Code)
	assert.Equal(t, "delivery agent is not available", strings.TrimSpace(assignRR2.Body.String()))
}
