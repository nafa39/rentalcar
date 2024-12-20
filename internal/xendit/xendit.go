package xendit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rental-car/internal/entity"
)

const apiUrl = "https://api.xendit.co/v2/invoices"

// CreateInvoice creates an invoice using Xendit API.
func CreateInvoice(product entity.ProductRequest, customer entity.CustomerRequest) (*entity.Invoice, error) {
	// Prepare the request body for Xendit API
	bodyRequest := map[string]interface{}{
		"external_id":      "1", // This could be dynamically set based on your needs (e.g., order ID, etc.)
		"amount":           product.Price,
		"description":      "Car Rental Invoice", // You can customize this based on the rental details
		"invoice_duration": 86400,                // 24 hours in seconds; adjust this if needed
		"customer": map[string]interface{}{
			"name":  customer.Name,
			"email": customer.Email,
		},
		"currency": "IDR", // Assuming the currency is IDR, change if necessary
		"items": []interface{}{
			map[string]interface{}{
				"name":     product.Name,
				"quantity": 1,
				"price":    product.Price,
			},
		},
	}

	// Marshal the body into JSON format
	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Prepare the HTTP client and the request
	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new HTTP request: %v", err)
	}

	// Get API Key from environment variable
	apiKey := os.Getenv("XENDIT_API_KEY") // Replace with your actual Xendit API key
	if apiKey == "" {
		return nil, fmt.Errorf("API Key is missing")
	}

	// Set the authentication and headers
	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	// Send the request to Xendit API
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer response.Body.Close()

	// Read the response body for debugging
	var responseBody []byte
	if responseBody, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Log the response for debugging
	fmt.Println("Response Body: ", string(responseBody))

	// Parse the response body into an Invoice object
	var resInvoice entity.Invoice
	if err := json.Unmarshal(responseBody, &resInvoice); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// Check if the invoice was created successfully
	if resInvoice.ID == "" || resInvoice.InvoiceURL == "" {
		return nil, fmt.Errorf("invoice ID or URL is empty")
	}

	// Return the created invoice
	return &resInvoice, nil
}
