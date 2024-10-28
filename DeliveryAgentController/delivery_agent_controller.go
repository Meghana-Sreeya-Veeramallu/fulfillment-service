package DeliveryAgentController

import (
	"Fulfillment/DeliveryAgentService"
	pb "Fulfillment/proto"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type DeliveryAgentServer struct {
	pb.UnimplementedDeliveryAgentServiceServer
	DB *gorm.DB
}

func NewDeliveryAgentServer(db *gorm.DB) *DeliveryAgentServer {
	return &DeliveryAgentServer{DB: db}
}

func (s *DeliveryAgentServer) AddDeliveryAgent(w http.ResponseWriter, r *http.Request) {
	var req pb.AddDeliveryAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := DeliveryAgentService.AddDeliveryAgent(s.DB, req.Name, req.City)
	if err != nil {
		if err.Error() == "name cannot be empty" || err.Error() == "city cannot be empty" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to add delivery agent: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pb.AddDeliveryAgentResponse{
		Message: "Delivery agent added successfully",
	})
}

func (s *DeliveryAgentServer) AssignAgentToOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID, err := strconv.Atoi(vars["agent-id"])
	if err != nil {
		http.Error(w, "Invalid agent ID", http.StatusBadRequest)
		return
	}

	orderID, err := strconv.Atoi(vars["order-id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	err = DeliveryAgentService.AssignAgentToOrder(s.DB, uint(agentID), orderID)
	if err != nil {
		if err.Error() == "delivery agent not found" || err.Error() == "delivery agent is not available" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "An error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&pb.AssignAgentToOrderResponse{
		Message: "Delivery agent assigned to order successfully",
	})
}

func (s *DeliveryAgentServer) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/delivery-agents", s.AddDeliveryAgent).Methods("POST")
	router.HandleFunc("/delivery-agents/{agent-id}/orders/{order-id}", s.AssignAgentToOrder).Methods("POST")
}
