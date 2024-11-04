package Service

import (
	client "Fulfillment/Client"
	"Fulfillment/Model"
	"Fulfillment/Repository"
	"errors"
)

// DeliveryAgentService handles business logic related to delivery agents.
type DeliveryAgentService struct {
	repo *Repository.DeliveryAgentRepository
}

// NewDeliveryAgentService creates a new DeliveryAgentService.
func NewDeliveryAgentService(repo *Repository.DeliveryAgentRepository) *DeliveryAgentService {
	return &DeliveryAgentService{repo: repo}
}

// AddDeliveryAgent adds a new delivery agent to the repository.
func (s *DeliveryAgentService) AddDeliveryAgent(name string, city string) (*Model.DeliveryAgent, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if city == "" {
		return nil, errors.New("city cannot be empty")
	}

	agent := &Model.DeliveryAgent{
		Name:               name,
		City:               city,
		AvailabilityStatus: Model.AVAILABLE,
		OrderID:            nil,
	}

	if err := s.repo.Create(agent); err != nil {
		return nil, err
	}

	return agent, nil
}

// AssignAgentToOrder assigns a delivery agent to an order.
func (s *DeliveryAgentService) AssignAgentToOrder(deliveryAgentID uint, orderID int) error {
	agent, err := s.repo.FindByID(deliveryAgentID)
	if err != nil {
		return errors.New("delivery agent not found")
	}

	if agent.AvailabilityStatus != Model.AVAILABLE {
		return errors.New("delivery agent is not available")
	}

	// Check if the order exists using the REST client function
	exists, err := client.CheckAndUpdateOrderStatus(orderID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("order does not exist")
	}

	agent.OrderID = &orderID
	agent.AvailabilityStatus = Model.UNAVAILABLE

	return s.repo.Save(agent)
}
