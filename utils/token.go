package utils

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyToken takes a token string and verifies it using the secret key
func VerifyToken(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Claims.(jwt.MapClaims); !ok {
			return "", fmt.Errorf("claim not supported")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if token != nil && token.Valid {
		return token.Claims.(jwt.MapClaims), nil
	}
	return nil, err
}

/*
GenerateToken generates a token from an input string
using the secret key and the HS256 signing method
*/
func GenerateToken(text string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(500 * time.Hour).Unix(),
		"sub": text,
	})
	tokStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokStr, nil
}

/*
GenerateTokenWithRole Generates a token with respect to the role of the admin
using an issued timestamp, expiration timestamp, subject, and role
*/
func GenerateTokenWithRole(text string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(500 * time.Hour).Unix(),
		"sub":  text,
		"role": role,
	})
	tokStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokStr, nil
}

// for account verification
func GenerateVerifyAccountToken() (string, error) {
	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func GeneratePIN() int {
	pin := rand.Intn(9000) + 1000
	// fmt.Println(pin)
	return pin
}
