package utils

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

func HashPassword(pass string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	return string(hashedByte), err
}

func ComparePassword(hashed, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
}

// Function to validate password strength
func IsValidPassword(password string) bool {
	// Check if the password has a minimum length
	if len(password) < 8 {
		return false
	}

	hasUpperCase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
			break
		}
	}

	hasLowerCase := false
	for _, char := range password {
		if unicode.IsLower(char) {
			hasLowerCase = true
			break
		}
	}

	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}

	// Check if the password meets all criteria
	return hasUpperCase && hasLowerCase && hasDigit
}

// OAuth2 configurations for different providers
var (
	GoogleConfig = &oauth2.Config{
		ClientID:     "your_google_client_id",
		ClientSecret: "your_google_client_secret",
		RedirectURL:  "your_redirect_url",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	FacebookConfig = &oauth2.Config{
		ClientID:     "your_facebook_client_id",
		ClientSecret: "your_facebook_client_secret",
		RedirectURL:  "your_redirect_url",
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}
	// appleConfig = &oauth2.Config{
	// 	ClientID:     "your_apple_client_id",
	// 	ClientSecret: "your_apple_client_secret",
	// 	RedirectURL:  "your_redirect_url",
	// 	Scopes:       []string{"email", "name"},
	// 	Endpoint:     apple.Endpoint,
	// }
)
