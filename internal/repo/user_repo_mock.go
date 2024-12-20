package repo

import (
	"github.com/stretchr/testify/mock"
)

// MockUserRepo is a mock implementation of the UserRepo interface
type MockUserRepo struct {
	mock.Mock
}

// UpdateBalance is the mocked method for updating the user's balance
func (m *MockUserRepo) UpdateBalance(userID int64, amount float64) error {
	args := m.Called(userID, amount)
	return args.Error(0)
}
