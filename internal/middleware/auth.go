package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware to protect routes
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
		}

		// Remove 'Bearer ' prefix if present
		token = strings.TrimPrefix(token, "Bearer ")

		fmt.Println("JWT Token received:", token) // Log the token for debugging

		userID, err := VerifyJWT(token)
		if err != nil {
			fmt.Println("Error verifying token:", err) // Log the error
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		fmt.Println("Token is valid. User ID:", userID) // Log the user ID from the token
		c.Set("userID", userID)
		return next(c)
	}
}
