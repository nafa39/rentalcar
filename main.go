package main

import (
	"log"
	"rental-car/config"
	handler "rental-car/internal/handlers"
	"rental-car/internal/repo"
	"rental-car/internal/validators"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDB()

	defer config.CloseDB(db)

	// Initialize Validator
	e.Validator = validators.NewCustomValidator()

	// Initialize repositories
	userRepo := repo.NewUserRepository(db)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userRepo)

	e.POST("/register", userHandler.RegisterUser)
	e.POST("/login", userHandler.LoginUser) // Login route

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
