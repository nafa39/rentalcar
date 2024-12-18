package main

import (
	"log"
	"rental-car/config"
	handler "rental-car/internal/handlers"
	"rental-car/internal/middleware"
	"rental-car/internal/repo"
	"rental-car/internal/validators"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := config.ConnectDB()

	defer config.CloseDB(db)

	// Add built-in middleware
	e.Use(echoMiddleware.Logger())  // Logs HTTP requests
	e.Use(echoMiddleware.Recover()) // Recovers from panics and logs them

	// Initialize Validator
	e.Validator = validators.NewCustomValidator()

	// Initialize repositories
	userRepo := repo.NewUserRepository(db)
	carRepo := repo.NewCarRepository(db)

	reservationRepo := repo.NewReservationRepository(db)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userRepo, reservationRepo)
	carHandler := handler.NewCarHandler(carRepo)

	e.POST("/register", userHandler.RegisterUser) // Register route
	e.POST("/login", userHandler.LoginUser)       // Login route

	// Protected routes
	secured := e.Group("/secure")
	secured.Use(middleware.JWTMiddleware)
	secured.POST("/top-up", userHandler.TopUpBalance)
	secured.POST("/rent", carHandler.RentCar) // Protected rent
	secured.GET("/booking/:reservationID", userHandler.GetBooking)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
