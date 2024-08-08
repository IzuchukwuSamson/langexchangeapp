package services

import "github.com/IzuchukwuSamson/lexi/internal/app/users/models"

type AdminServiceInterface interface {
	FetchAllUsers() ([]models.User, error)
	AdminEmailExists(email string) (bool error)
	NewAdminEmail()
}
