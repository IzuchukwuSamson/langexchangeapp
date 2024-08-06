package utils

import (
	"time"

	"github.com/IzuchukwuSamson/lexi/internal/app/users/models"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(user *models.User) (string, error) {
	signingKey := []byte("your_secret_key") // Use a secure secret key

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	})

	// Sign the token
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
