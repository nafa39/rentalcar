package handler

import (
	"log"
	"net/http"
	"rental-car/internal/entity"
	"rental-car/internal/repo"
	"rental-car/internal/xendit"
	"time"

	"github.com/labstack/echo/v4"
)

// CarHandler contains dependencies for car-related handlers.
type CarHandler struct {
	CarRepo repo.CarRepository
}

// NewCarHandler creates a new CarHandler instance.
func NewCarHandler(carRepo repo.CarRepository) *CarHandler {
	return &CarHandler{CarRepo: carRepo}
}

type XenditService interface {
	CreateInvoice(product entity.ProductRequest, customer entity.CustomerRequest) (*entity.Invoice, error)
}

// RentCar handles the rent car request.
func (h *CarHandler) RentCar(c echo.Context) error {
	// Get userID from the JWT middleware
	userID := c.Get("userID").(int64)

	// Parse the request body
	var req struct {
		CarID     int64  `json:"car_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Define the date formats
	layouts := []string{
		"2006-01-02T15:04:05Z07:00", // Full datetime
		"2006-01-02",                // Date-only format
	}

	// Parse start and end dates
	var startDate, endDate time.Time
	var err error
	for _, layout := range layouts {
		startDate, err = time.Parse(layout, req.StartDate)
		if err == nil {
			break
		}
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start date format"})
	}

	for _, layout := range layouts {
		endDate, err = time.Parse(layout, req.EndDate)
		if err == nil {
			break
		}
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end date format"})
	}

	// Call the repository to rent the car
	reservationID, totalPrice, err := h.CarRepo.RentCar(userID, req.CarID, startDate, endDate)
	if err != nil {
		log.Printf("Error renting car: %v", err)
		if err.Error() == "insufficient balance" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient balance"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error renting car"})
	}

	// Fetch car details (for invoice creation)
	car, err := h.CarRepo.GetCarByID(req.CarID)
	if err != nil {
		log.Printf("Error fetching car details: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching car details"})
	}

	// Fetch user details (for invoice creation)
	user, err := h.CarRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Error fetching user details: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error fetching user details"})
	}

	// Create an invoice using Xendit API
	invoice, err := xendit.CreateInvoice(entity.ProductRequest{
		Name:  car.Name,
		Price: totalPrice,
	}, entity.CustomerRequest{
		Name:  user.Name,
		Email: user.Email,
	})
	if err != nil {
		log.Printf("Error creating invoice: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating invoice"})
	}

	// Respond with the rental and invoice details
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "Car rented successfully",
		"reservation_id": reservationID,
		"total_price":    totalPrice,
		"invoice_id":     invoice.ID,
		"invoice_url":    invoice.InvoiceURL,
	})
}
