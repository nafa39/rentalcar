package handler

import (
	"net/http"
	"rental-car/internal/entity"
	"rental-car/internal/repo"
	"rental-car/internal/utils"

	"github.com/labstack/echo/v4"
)

// UserHandler contains dependencies for user-related handlers.
type UserHandler struct {
	UserRepo repo.UserRepository
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userRepo repo.UserRepository) *UserHandler {
	return &UserHandler{UserRepo: userRepo}
}

// RegisterUser handles the user registration request.
func (h *UserHandler) RegisterUser(c echo.Context) error {
	// Parse request body into a User entity
	var req entity.User
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	// Validate input
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not hash password"})
	}

	// Create a new user entity
	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword, // Hash the password
	}

	// Save user in the database
	userID, err := h.UserRepo.Create(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not register user"})
	}

	// Respond with the newly created user's ID
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user_id": userID,
	})
}
