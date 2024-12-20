package repo

import (
	"rental-car/internal/entity"

	"github.com/stretchr/testify/mock"
)

// MockReservationRepository is a mock implementation of the ReservationRepository interface.
type MockReservationRepository struct {
	mock.Mock
}

// GetBooking mocks the GetBooking method of the ReservationRepository interface.
func (m *MockReservationRepository) GetBooking(userID, reservationID int64) (*entity.Reservation, error) {
	args := m.Called(userID, reservationID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Reservation), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetAllBookings mocks the GetAllBookings method of the ReservationRepository interface.
func (m *MockReservationRepository) GetAllBookings(userID int64) ([]entity.Reservation, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]entity.Reservation), args.Error(1)
	}
	return nil, args.Error(1)
}
