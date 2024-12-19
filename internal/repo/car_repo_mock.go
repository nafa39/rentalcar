package repo

import (
	"fmt"
	"rental-car/internal/entity"
	"time"

	"github.com/stretchr/testify/mock"
)

type CarRepoMock struct {
	mock.Mock
}

func (m *CarRepoMock) RentCar(userID, carID int64, startDate, endDate time.Time) (int64, float64, error) {
	// Mock the check for car availability
	carArgs := m.Mock.Called(carID)
	if carArgs.Get(0) == nil {
		return 0, 0, fmt.Errorf("car not available")
	}
	car := carArgs.Get(0).(entity.Car)

	// Mock the check for user existence
	userArgs := m.Mock.Called(userID)
	if userArgs.Get(0) == nil {
		return 0, 0, fmt.Errorf("user not found")
	}
	user := userArgs.Get(0).(entity.User)

	// Calculate duration and total price
	duration := endDate.Sub(startDate).Hours() / 24
	totalPrice := car.PricePerDay * duration

	// Check if user has sufficient balance
	if user.Balance < totalPrice {
		return 0, 0, fmt.Errorf("insufficient balance")
	}

	// Return the mock reservation ID and total price
	return 123, totalPrice, nil
}
