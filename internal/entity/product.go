package entity

type ProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Invoice struct {
	ID          string  `json:"id"`
	ExternalID  string  `json:"external_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	InvoiceURL  string  `json:"invoice_url"`
}
