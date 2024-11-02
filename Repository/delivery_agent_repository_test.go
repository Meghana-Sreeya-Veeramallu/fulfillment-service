package Repository

import (
	"Fulfillment/Model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// setupDB function to initialize a new in-memory database for testing.
func setupDB() (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Model.DeliveryAgent{})
	if err != nil {
		return nil, nil
	}
	return db, func() {
		db.Exec("DROP TABLE delivery_agents")
	}
}

// Test Create
func TestCreateDeliveryAgent(t *testing.T) {
	db, teardown := setupDB()
	defer teardown()

	repo := NewDeliveryAgentRepository(db)
	agent := &Model.DeliveryAgent{Name: "Test Agent"}

	if err := repo.Create(agent); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if agent.Id == 0 {
		t.Fatalf("expected agent ID to be set")
	}
}

// Test FindByID
func TestFindByID(t *testing.T) {
	db, teardown := setupDB()
	defer teardown()

	repo := NewDeliveryAgentRepository(db)
	agent := &Model.DeliveryAgent{Name: "Test Agent"}
	err := repo.Create(agent)
	if err != nil {
		return
	}

	foundAgent, err := repo.FindByID(agent.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if foundAgent.Id != agent.Id {
		t.Fatalf("expected agent ID %v, got %v", agent.Id, foundAgent.Id)
	}
}

// Test Save
func TestSaveDeliveryAgent(t *testing.T) {
	db, teardown := setupDB()
	defer teardown()

	repo := NewDeliveryAgentRepository(db)
	agent := &Model.DeliveryAgent{Name: "Test Agent"}
	err := repo.Create(agent)
	if err != nil {
		return
	}

	agent.Name = "Updated Agent"
	if err := repo.Save(agent); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedAgent, err := repo.FindByID(agent.Id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedAgent.Name != "Updated Agent" {
		t.Fatalf("expected agent name to be 'Updated Agent', got %v", updatedAgent.Name)
	}
}
