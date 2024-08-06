package utils

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/IzuchukwuSamson/lexi/internal/app/users/models"
)

func ToUserDTO(user models.User) models.UserDTO {
	return models.UserDTO{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}
}

func SanitizeInput(s string) string {
	// Trim leading and trailing whitespace
	s = strings.TrimSpace(s)

	// Remove extra spaces within the string
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	if len(s) == 0 {
		return s
	}

	// Handle Unicode characters and capitalize the first letter
	r, size := utf8.DecodeRuneInString(s)
	if !unicode.IsLetter(r) {
		return s
	}

	// Convert first letter to uppercase and the rest to lowercase
	first := string(unicode.ToUpper(r))
	rest := strings.ToLower(s[size:])
	return first + rest
}
