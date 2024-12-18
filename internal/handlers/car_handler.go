package handler

import (
	"log"
	"net/http"
	"rental-car/internal/repo"
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

// RentCar handles the rent car request.
func (h *CarHandler) RentCar(c echo.Context) error {
	// Get userID from the JWT middleware
	userID := c.Get("userID").(int64)

	// Parse the request body
	var req struct {
		CarID     int64     `json:"car_id"`
		StartDate time.Time `json:"start_date"`
		EndDate   time.Time `json:"end_date"`
	}

	if err := c.Bind(&req); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Call the repository to rent the car
	reservationID, totalPrice, err := h.CarRepo.RentCar(userID, req.CarID, req.StartDate, req.EndDate)
	if err != nil {
		log.Printf("Error renting car: %v", err)
		if err.Error() == "insufficient balance" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Insufficient balance"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error renting car"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "Car rented successfully",
		"reservation_id": reservationID,
		"total_price":    totalPrice,
	})
}
