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

// CreateInvoice creates an invoice using Xendit API
func CreateInvoice(product entity.ProductRequest, customer entity.CustomerRequest) (*entity.Invoice, error) {
	bodyRequest := map[string]interface{}{
		"external_id":      "1",
		"amount":           product.Price,
		"description":      "Dummy Invoice RMT003",
		"invoice_duration": 86400,
		"customer": map[string]interface{}{
			"name":  customer.Name,
			"email": customer.Email,
		},
		"currency": "IDR",
		"items": []interface{}{
			map[string]interface{}{
				"name":     product.Name,
				"quantity": 1,
				"price":    product.Price,
			},
		},
	}

	reqBody, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	var apiKey = os.Getenv("XENDIT_API_KEY") // Replace with your actual Xendit API key
	if apiKey == "" {
		return nil, fmt.Errorf("API Key is missing!")
	}

	request.SetBasicAuth(apiKey, "")
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Log the raw response body for debugging
	var responseBody []byte
	if responseBody, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	// Print the response body to debug
	fmt.Println("Response Body: ", string(responseBody))

	// Decode the response body
	var resInvoice entity.Invoice
	if err := json.Unmarshal(responseBody, &resInvoice); err != nil {
		return nil, err
	}

	// Check if Invoice fields are set
	if resInvoice.ID == "" || resInvoice.InvoiceURL == "" {
		return nil, fmt.Errorf("invoice ID or URL is empty")
	}

	return &resInvoice, nil
}
