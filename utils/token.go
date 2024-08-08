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
// func GenerateToken(text string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"iat": time.Now().Unix(),
// 		"exp": time.Now().Add(500 * time.Hour).Unix(),
// 		"sub": text,
// 	})
// 	tokStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokStr, nil
// }

func GenerateToken(text string) (string, string, error) {
	// Generate Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
		"sub": text,
	})
	accessTokStr, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	// Generate Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days expiration
		"sub": text,
	})
	refreshTokStr, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}

	return accessTokStr, refreshTokStr, nil
}

func RefreshAccessToken(refreshTokenStr string) (string, error) {
	// Parse the refresh token
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return "", err
	}

	// Validate the token and extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Generate new access token
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(1 * time.Hour).Unix(), // 1 hour expiration
			"sub": claims["sub"],
		})
		newAccessTokStr, err := newAccessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			return "", err
		}
		return newAccessTokStr, nil
	} else {
		return "", err
	}
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
