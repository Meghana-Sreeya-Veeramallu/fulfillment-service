package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test CheckOrderExists when the order exists
func TestCheckOrderExistsWhenOrderExists(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Order exists
	}))
	defer server.Close()

	checkOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	exists, err := CheckOrderExists(1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("expected order to exist, got false")
	}
}

// Test CheckOrderExists when the order does not exist
func TestCheckOrderExistsWhenOrderNotFound(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound) // Order does not exist
	}))
	defer server.Close()

	// Mock the checkOrderExists function
	checkOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	exists, err := CheckOrderExists(2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if exists {
		t.Errorf("expected order to not exist, got true")
	}
}

// Test CheckOrderExists when an unexpected error occurs
func TestCheckOrderExists_UnexpectedError(t *testing.T) {
	// Setup mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError) // Simulate an unexpected error
	}))
	defer server.Close()

	// Mock the checkOrderExists function
	checkOrderExists = func(orderID int) (*http.Response, error) {
		return http.Get(fmt.Sprintf("%s/orders/%d", server.URL, orderID))
	}

	exists, err := CheckOrderExists(3)
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
