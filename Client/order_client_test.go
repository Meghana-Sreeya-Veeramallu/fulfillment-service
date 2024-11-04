package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test CheckAndUpdateOrderStatus when the order exists
func TestCheckAndUpdateOrderStatusWhenOrderExists(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Simulate order exists
	}))
	defer server.Close()

	CheckOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	// Mock updateOrderStatus to simulate a successful update
	UpdateOrderStatus = func(orderID int) error {
		return nil
	}

	exists, err := CheckAndUpdateOrderStatus(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("expected order to exist, got false")
	}
}

// TestCheckAndUpdateOrderStatus tests the case where the order cannot be assigned
func TestCheckAndUpdateOrderStatusWhenOrderCannotBeAssigned(t *testing.T) {
	// Setup mock server to simulate order exists
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Simulate order exists
	}))
	defer server.Close()

	// Mock CheckAndUpdateOrderStatus to use the mock server
	CheckOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	// Mock UpdateOrderStatus to simulate the assignment failure
	UpdateOrderStatus = func(orderID int) error {
		return fmt.Errorf("order cannot be assigned") // Simulate order cannot be assigned
	}

	exists, err := CheckAndUpdateOrderStatus(1)
	if err == nil || err.Error() != "order cannot be assigned" {
		t.Errorf("expected order cannot be assigned error, got %v", err)
	}
	if !exists {
		t.Errorf("expected order to exist, got false")
	}
}

func TestCheckAndUpdateOrderStatusWhenOrderExistsUpdateFails(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Simulate order exists
	}))
	defer server.Close()

	CheckOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	// Mock updateOrderStatus to simulate a failure
	UpdateOrderStatus = func(orderID int) error {
		return fmt.Errorf("update failed")
	}

	exists, err := CheckAndUpdateOrderStatus(1)
	if err == nil || err.Error() != "update failed" {
		t.Errorf("expected update failed error, got %v", err)
	}
	if !exists {
		t.Errorf("expected order to exist, got false")
	}
}

// Test CheckAndUpdateOrderStatus when the order does not exist
func TestCheckAndUpdateOrderStatusWhenOrderNotFound(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // Order does not exist
	}))
	defer server.Close()

	// Mock the checkOrderExists function
	CheckOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	exists, err := CheckAndUpdateOrderStatus(2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if exists {
		t.Errorf("expected order to not exist, got true")
	}
}

// Test CheckAndUpdateOrderStatus when an unexpected error occurs
func TestCheckAndUpdateOrderStatusWhenUnexpectedError(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Simulate an unexpected error
	}))
	defer server.Close()

	// Mock the checkOrderExists function
	CheckOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	exists, err := CheckAndUpdateOrderStatus(3)
	if err == nil {
		t.Errorf("expected an error, got none")
	}
	if exists {
		t.Errorf("expected order to not exist, got true")
	}

	// Check if the error message indicates a failure
	expectedErrMsg := "unexpected status code: 500"
	if !strings.Contains(err.Error(), expectedErrMsg) {
		t.Errorf("expected error to contain %q, got %q", expectedErrMsg, err.Error())
	}
}
