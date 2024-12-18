package repo

import (
	"rental-car/internal/entity"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// UserRepository defines the methods for interacting with the User entity.
type UserRepository interface {
	Create(user *entity.User) (int64, error) // Create a new user and return the ID
	// RegisterUser(user *entity.User) (int64, error)        // Register a new user
	// UpdateBalance(userID int64, newBalance float64) error // Update the user's balance
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Create a new user and return the ID
func (r *userRepo) Create(user *entity.User) (int64, error) {
	if err := r.db.Create(user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}
