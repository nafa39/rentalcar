package middleware

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// Create JWT Token
func CreateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token using the claims and a secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
