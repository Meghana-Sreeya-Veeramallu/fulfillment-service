package Model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	var err error
	dsn := "host=localhost user=postgres password=admin dbname=fulfillment_service port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	if err := db.AutoMigrate(&DeliveryAgent{}); err != nil {
		panic("failed to migrate database schema")
	}

	tx := db.Begin()

	code := m.Run()

	if err := db.Migrator().DropTable(&DeliveryAgent{}); err != nil {
		panic("failed to drop table")
	}

	tx.Rollback()

	os.Exit(code)
}

func TestNewDeliveryAgent(t *testing.T) {
	name := "Ketan"
	city := "Hyderabad"

	agent, err := NewDeliveryAgent(db, name, city)
	if err != nil {
		t.Fatalf("Failed to create delivery agent: %v", err)
	}

	if agent.Id == 0 {
		t.Error("Expected auto-incremented ID to be set, got 0")
	}

	if agent.Name != name {
		t.Errorf("Expected Name to be %s, got %s", name, agent.Name)
	}

	if agent.City != city {
		t.Errorf("Expected City to be %s, got %s", city, agent.City)
	}

	if agent.AvailabilityStatus != AVAILABLE {
		t.Errorf("Expected AvailabilityStatus to be %s, got %s", AVAILABLE, agent.AvailabilityStatus)
	}

	if agent.OrderID != nil {
		t.Error("Expected OrderID to be nil")
	}
}

func TestNewDeliveryAgentWithEmptyName(t *testing.T) {
	agent, err := NewDeliveryAgent(db, "", "Hyderabad")
	if err == nil {
		t.Fatal("Expected error for empty name, got none")
	}
	if agent != nil {
		t.Error("Expected agent to be nil")
	}
}

func TestNewDeliveryAgentWithEmptyCity(t *testing.T) {
	agent, err := NewDeliveryAgent(db, "Ketan", "")
	if err == nil {
		t.Fatal("Expected error for empty city, got none")
	}
	if agent != nil {
		t.Error("Expected agent to be nil")
	}
}

func TestNewDeliveryAgentWithEmptyNameAndCity(t *testing.T) {
	agent, err := NewDeliveryAgent(db, "", "")
	if err == nil {
		t.Fatal("Expected error for empty name and city, got none")
	}
	if agent != nil {
		t.Error("Expected agent to be nil")
	}
}
