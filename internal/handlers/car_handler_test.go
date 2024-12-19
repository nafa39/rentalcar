package handler

import (
	"fmt"
	"rental-car/internal/entity"
	"strings"
	"testing"
	"time"

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCarRepo is a mock of the CarRepository interface.
type MockCarRepo struct {
	mock.Mock
}

func (m *MockCarRepo) RentCar(userID, carID int64, startDate, endDate time.Time) (int64, float64, error) {
	args := m.Called(userID, carID, startDate, endDate)
	return args.Get(0).(int64), args.Get(1).(float64), args.Error(2)
}

func TestRentCar_Success(t *testing.T) {
	// Create a new instance of Echo
	e := echo.New()

	// Create a mock of the CarRepo
	mockRepo := new(MockCarRepo)

	// Create the handler with the mock repo
	handler := NewCarHandler(mockRepo)

	// Define the request body
	reqBody := map[string]interface{}{
		"car_id":     101,
		"start_date": "2024-12-20T02:24:46Z",
		"end_date":   "2024-12-21T02:24:46Z",
	}

	// Marshal the request body to JSON
	reqJSON, _ := json.Marshal(reqBody)

	// Create a new HTTP request with the JSON body
	req := httptest.NewRequest(http.MethodPost, "/rent-car", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// Set up the user ID in the context (this should be mocked as well)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("userID", int64(1)) // Mocking the user ID from JWT middleware

	// Mock the RentCar function to return the expected values
	mockRepo.On("RentCar", int64(1), int64(101), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(int64(12345), 200.0, nil)

	// Call the RentCar handler
	err := handler.RentCar(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body
	expectedResp := `{"message":"Car rented successfully","reservation_id":12345,"total_price":200}`
	assert.JSONEq(t, expectedResp, rec.Body.String())

	// Assert that the mock's RentCar method was called exactly once with the expected arguments
	mockRepo.AssertExpectations(t)
}

func TestRentCar_InsufficientBalance(t *testing.T) {
	// Prepare test data
	carID := int64(101)
	userID := int64(1)
	startDate := time.Now().Truncate(time.Second)                  // Truncate to the nearest second
	endDate := startDate.Add(24 * time.Hour).Truncate(time.Second) // Truncate to the nearest second

	// Create a mock of the CarRepo
	mockRepo := new(MockCarRepo)

	// Create the handler with the mock repo
	handler := NewCarHandler(mockRepo)

	// Mock RentCar repository method to simulate insufficient balance scenario
	mockRepo.On("RentCar", userID, carID, startDate, endDate).Return(int64(0), 0.0, fmt.Errorf("insufficient balance"))

	// Mock car availability
	car := entity.Car{
		ID:          carID,
		PricePerDay: 100.0, // PricePerDay set to 100
		Status:      "available",
	}
	mockRepo.On("RentCar", carID).Return(car, nil)

	// Mock user with insufficient balance
	user := entity.User{
		ID:      userID,
		Balance: 50.0, // Insufficient balance
	}
	mockRepo.On("RentCar", userID).Return(user, nil)

	// Create Echo instance for testing
	e := echo.New()

	// Prepare the request with car_id, start_date, and end_date
	req := fmt.Sprintf(`{"car_id": %d, "start_date": "%s", "end_date": "%s"}`,
		carID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))
	rec := httptest.NewRecorder()
	reqReader := httptest.NewRequest("POST", "/rent", strings.NewReader(req))

	// Set Content-Type header to application/json
	reqReader.Header.Set("Content-Type", "application/json")

	// Create a new context with the request and response recorder
	c := e.NewContext(reqReader, rec)

	// Set the userID in context (simulating logged-in user)
	c.Set("userID", userID)

	// Call the RentCar handler
	err := handler.RentCar(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)              // Expecting 400 status code
	assert.Contains(t, rec.Body.String(), "Insufficient balance") // Error message check
}

func TestRentCar_CarNotAvailable(t *testing.T) {
	// Prepare test data
	carID := int64(101)
	userID := int64(1)
	startDate := time.Now().Truncate(time.Second)                  // Truncate to the nearest second
	endDate := startDate.Add(24 * time.Hour).Truncate(time.Second) // Truncate to the nearest second

	// Create a mock of the CarRepo
	mockRepo := new(MockCarRepo)

	// Create the handler with the mock repo
	handler := NewCarHandler(mockRepo)

	// Mock RentCar repository method to simulate car not available scenario
	mockRepo.On("RentCar", userID, carID, startDate, endDate).Return(int64(0), 0.0, fmt.Errorf("car not available"))

	// Prepare Echo instance for testing
	e := echo.New()

	// Prepare the request with car_id, start_date, and end_date
	req := fmt.Sprintf(`{"car_id": %d, "start_date": "%s", "end_date": "%s"}`,
		carID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))
	rec := httptest.NewRecorder()
	reqReader := httptest.NewRequest("POST", "/rent", strings.NewReader(req))

	// Set Content-Type header to application/json
	reqReader.Header.Set("Content-Type", "application/json")

	// Create a new context with the request and response recorder
	c := e.NewContext(reqReader, rec)

	// Set the userID in context (simulating logged-in user)
	c.Set("userID", userID)

	// Call the RentCar handler
	err := handler.RentCar(c)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)  // Expecting 500 status code
	assert.Contains(t, rec.Body.String(), "Error renting car") // Error message check
}
