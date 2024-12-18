package handler

import (
	"log"
	"net/http"
	"rental-car/internal/entity"
	"rental-car/internal/middleware"
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
	var user entity.User
	if err := c.Bind(&user); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	// Validate input
	if err := c.Validate(user); err != nil {
		log.Printf("Error validating request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if user.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password is required"})
	}

	// Hash the user's password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error hashing password"})
	}
	//log.Println("Hashed password:", hashedPassword)
	user.Password = hashedPassword

	//log.Println("User password:", user.Password)

	userID, err := h.UserRepo.RegisterUser(&user)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating user"})
	}

	// Respond with the newly created user's ID
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "user registered successfully",
		"user_id": userID,
	})
}

// LoginUser handles user login
func (h *UserHandler) LoginUser(c echo.Context) error {
	var loginData struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&loginData); err != nil {
		log.Printf("Error binding request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Find user by email
	user, err := h.UserRepo.FindByEmail(loginData.Email)
	if err != nil {
		log.Printf("User not found: %v", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Check if the password is correct
	if !utils.ComparePassword(loginData.Password, user.Password) {
		log.Println("User password:", user.Password)
		log.Println("Invalid password")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Create JWT token
	token, err := middleware.CreateJWT(user.ID)
	if err != nil {
		log.Printf("Error creating JWT: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error creating token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
