package repo

import "rental-car/internal/entity"

// ReservationRepository defines the methods for interacting with the Reservation entity.
type ReservationRepository interface {
	Create(reservation *entity.Reservation) (int64, error) // Create a new reservation and return the ID
}
