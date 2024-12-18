package repo

import (
	"fmt"
	"rental-car/internal/entity"
	"time"

	"gorm.io/gorm"
)

// CarRepository defines the methods for handling car-related database operations.
type CarRepository interface {
	RentCar(userID, carID int64, startDate, endDate time.Time) (int64, float64, error)
}

// carRepo implements CarRepository interface
type carRepo struct {
	db *gorm.DB
}

// NewCarRepository creates a new car repository.
func NewCarRepository(db *gorm.DB) CarRepository {
	return &carRepo{db: db}
}

// RentCar allows a user to rent a car for a specific period.
func (r *carRepo) RentCar(userID, carID int64, startDate, endDate time.Time) (int64, float64, error) {
	// Check if car is available
	var car entity.Car
	if err := r.db.Where("id = ? AND status = ?", carID, "available").First(&car).Error; err != nil {
		return 0, 0, fmt.Errorf("car not available")
	}

	// Check if the user exists
	var user entity.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return 0, 0, fmt.Errorf("user not found")
	}

	// Calculate the total price for the rental (price_per_day * number of days)
	duration := endDate.Sub(startDate).Hours() / 24 // Duration in days
	totalPrice := car.PricePerDay * duration

	// Check if the user has sufficient balance
	if user.Balance < totalPrice {
		return 0, 0, fmt.Errorf("insufficient balance")
	}

	// Create the reservation
	reservation := entity.Reservation{
		UserID:     userID,
		CarID:      carID,
		StartDate:  startDate,
		EndDate:    endDate,
		TotalPrice: totalPrice,
	}

	if err := r.db.Create(&reservation).Error; err != nil {
		return 0, 0, err
	}

	// Deduct balance from the user
	user.Balance -= totalPrice
	if err := r.db.Save(&user).Error; err != nil {
		return 0, 0, err
	}

	// Mark the car as rented
	car.Status = "rented"
	if err := r.db.Save(&car).Error; err != nil {
		return 0, 0, err
	}

	return reservation.ID, totalPrice, nil
}
