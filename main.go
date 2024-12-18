package main

import (
	"log"
	"rental-car/config"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
