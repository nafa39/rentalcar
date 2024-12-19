package repo

import (
	"fmt"
	"rental-car/internal/entity"

	"gorm.io/gorm"
)

// ReservationRepository defines the methods for interacting with the Reservation entity.
type ReservationRepository interface {
	GetBooking(userID, reservationID int64) (*entity.Reservation, error) // Get booking details by userID and reservationID
	GetAllBookings(userID int64) ([]entity.Reservation, error)           // Get all bookings by userID
}

// reservationRepo implements the ReservationRepository interface
type reservationRepo struct {
	db *gorm.DB
}

// NewReservationRepository creates a new instance of ReservationRepository
func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepo{db: db}
}

// GetBooking fetches a booking by userID and reservationID
func (r *reservationRepo) GetBooking(userID, reservationID int64) (*entity.Reservation, error) {
	var reservation entity.Reservation

	// Query to get the booking where the user_id matches the user and reservation_id matches
	if err := r.db.Preload("Car").Where("user_id = ? AND id = ?", userID, reservationID).First(&reservation).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, err
	}

	return &reservation, nil
}

// GetAllBookings fetches all bookings for a user
func (r *reservationRepo) GetAllBookings(userID int64) ([]entity.Reservation, error) {
	var reservations []entity.Reservation

	// Query to get all bookings where the user_id matches
	if err := r.db.Preload("Car").Where("user_id = ?", userID).Find(&reservations).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no bookings found for this user")
		}
		return nil, err
	}

	return reservations, nil
}
