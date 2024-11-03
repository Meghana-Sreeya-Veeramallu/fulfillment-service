package client

import (
	"fmt"
	"net/http"
)

// checkOrderExists makes the actual HTTP call.
var checkOrderExists = func(orderID int) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8081/orders/%d", orderID)
	return http.Get(url)
}

// CheckOrderExists checks if the order exists by making an HTTP GET request.
func CheckOrderExists(orderID int) (bool, error) {
	resp, err := checkOrderExists(orderID)
	if err != nil {
		return false, fmt.Errorf("failed to call order service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil // Order exists
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil // Order does not exist
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
