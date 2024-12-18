package repo

import (
	"errors"
	"rental-car/internal/entity"

	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// UserRepository defines the methods for interacting with the User entity.
type UserRepository interface {
	Create(user *entity.User) (int64, error)          // Create a new user and return the ID
	FindByEmail(email string) (*entity.User, error)   // Add FindByEmail method
	RegisterUser(user *entity.User) (int64, error)    // Register a new user
	UpdateBalance(userID int64, amount float64) error // Update the user's balance
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

// FindByEmail retrieves a user by their email address
func (r *userRepo) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// RegisterUser creates a new user (this could be an alias of Create)
func (r *userRepo) RegisterUser(user *entity.User) (int64, error) {
	return r.Create(user)
}

// UpdateBalance updates the user's balance
func (r *userRepo) UpdateBalance(userID int64, amount float64) error {
	var user entity.User

	// Fetch the user
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Update the balance
	newBalance := user.Balance + amount
	if newBalance < 0 {
		return errors.New("balance cannot be negative")
	}

	user.Balance = newBalance
	return r.db.Save(&user).Error
}
