package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rental-car/internal/entity"
	"rental-car/internal/repo"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TestGetBookingSuccess tests the successful retrieval of bookings.
func TestGetBookingSuccess(t *testing.T) {
	e := echo.New()
	mockRepo := new(repo.MockReservationRepository)
	handler := &UserHandler{ReservationRepo: mockRepo}

	// Mock data
	userID := int64(1)
	mockBookings := []entity.Reservation{
		{
			Car:        entity.Car{Name: "Toyota Camry", Category: "Sedan"},
			StartDate:  time.Date(2024, 12, 20, 14, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2024, 12, 22, 14, 0, 0, 0, time.UTC),
			TotalPrice: 300.00,
		},
		{
			Car:        entity.Car{Name: "Honda Civic", Category: "Compact"},
			StartDate:  time.Date(2024, 12, 25, 10, 0, 0, 0, time.UTC),
			EndDate:    time.Date(2024, 12, 27, 10, 0, 0, 0, time.UTC),
			TotalPrice: 200.00,
		},
	}

	// Mock the repository call
	mockRepo.On("GetAllBookings", userID).Return(mockBookings, nil)

	// Create a new request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	// Call the handler
	err := handler.GetBooking(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Verify the response matches the expected format
	expectedResponse := []map[string]interface{}{
		{
			"car_name":    "Toyota Camry",
			"category":    "Sedan",
			"start_date":  "2024-12-20T14:00:00Z",
			"end_date":    "2024-12-22T14:00:00Z",
			"total_price": 300.00,
		},
		{
			"car_name":    "Honda Civic",
			"category":    "Compact",
			"start_date":  "2024-12-25T10:00:00Z",
			"end_date":    "2024-12-27T10:00:00Z",
			"total_price": 200.00,
		},
	}

	var actualResponse []map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)

	// Ensure the mock method was called as expected
	mockRepo.AssertExpectations(t)
}

// TestGetBookingErrorFetching tests the error scenario when fetching bookings fails.
func TestGetBookingErrorFetching(t *testing.T) {
	e := echo.New()
	mockRepo := new(repo.MockReservationRepository)
	handler := &UserHandler{ReservationRepo: mockRepo}

	// Mock data
	userID := int64(1)
	mockError := fmt.Errorf("failed to fetch bookings")

	// Mock the repository call to return an error
	mockRepo.On("GetAllBookings", userID).Return(nil, mockError)

	// Create a new request and response recorder
	req := httptest.NewRequest(http.MethodGet, "/bookings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", userID)

	// Call the handler
	err := handler.GetBooking(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Verify the response contains the error message
	expectedResponse := map[string]string{
		"error": mockError.Error(),
	}

	var actualResponse map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)

	// Ensure the mock method was called as expected
	mockRepo.AssertExpectations(t)
}
