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

	// // Set your secret API key from Xendit
	// apiKey := os.Getenv("XENDIT_API_KEY")
	// if apiKey == "" {
	// 	log.Fatal("API Key is missing!")
	// }

	// // Example inputs for testing
	// product := entity.ProductRequest{
	// 	Name:  "Sample Product",
	// 	Price: 500000.00,
	// }
	// customer := entity.CustomerRequest{
	// 	Name:  "John Doe",
	// 	Email: "john@example.com",
	// }

	// // Create the invoice using Xendit
	// invoice, err := xendit.CreateInvoice(product, customer)
	// if err != nil {
	// 	log.Fatalf("Error creating invoice: %v", err)
	// }

	// // Print the invoice URL to check if it was generated
	// fmt.Println("Invoice successfully created!")
	// fmt.Printf("Invoice ID: %s\n", invoice.ID)
	// fmt.Printf("Invoice URL: %s\n", invoice.InvoiceURL)

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
	secured.GET("/booking", userHandler.GetBooking)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
