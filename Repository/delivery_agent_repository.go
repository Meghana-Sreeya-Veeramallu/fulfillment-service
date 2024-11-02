package Repository

import (
	"Fulfillment/Model"
	"gorm.io/gorm"
)

// DeliveryAgentRepository provides methods to interact with the database for Delivery Agents.
type DeliveryAgentRepository struct {
	DB *gorm.DB
}

// NewDeliveryAgentRepository creates a new DeliveryAgentRepository.
func NewDeliveryAgentRepository(db *gorm.DB) *DeliveryAgentRepository {
	return &DeliveryAgentRepository{DB: db}
}

// Create adds a new delivery agent to the database.
func (r *DeliveryAgentRepository) Create(agent *Model.DeliveryAgent) error {
	return r.DB.Create(agent).Error
}

// FindByID retrieves a delivery agent by ID.
func (r *DeliveryAgentRepository) FindByID(id uint) (*Model.DeliveryAgent, error) {
	var agent Model.DeliveryAgent
	if err := r.DB.First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// Save updates a delivery agent in the database.
func (r *DeliveryAgentRepository) Save(agent *Model.DeliveryAgent) error {
	return r.DB.Save(agent).Error
}
