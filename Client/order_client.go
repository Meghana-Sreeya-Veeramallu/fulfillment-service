package client

import (
	"bytes"
	"fmt"
	"net/http"
)

// CheckOrderExists makes the actual HTTP call.
var CheckOrderExists = func(orderID int) (*http.Response, error) {
	url := fmt.Sprintf("http://localhost:8081/orders/%d", orderID)
	return http.Get(url)
}

// UpdateOrderStatus makes an HTTP PUT request to update the order status.
var UpdateOrderStatus = func(orderID int) error {
	url := fmt.Sprintf("http://localhost:8081/orders/%d/status", orderID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(nil))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Order status updated successfully for order ID: %d\n", orderID)
		return nil
	}
	if resp.StatusCode == http.StatusBadRequest {
		return fmt.Errorf("order cannot be assigned")
	}
	return fmt.Errorf("unexpected status code when updating order status: %d", resp.StatusCode)
}

// CheckAndUpdateOrderStatus checks if an order exists and updates its status.
var CheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
	resp, err := CheckOrderExists(orderID)
	if err != nil {
		return false, fmt.Errorf("failed to call order service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if err := UpdateOrderStatus(orderID); err != nil {
			return true, err // Order exists but the update fails
		}
		return true, nil // Order exists and was updated
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil // Order does not exist
	}

	return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
