package services

import (
	adminModel "github.com/IzuchukwuSamson/lexi/internal/app/admin/models"
	userModel "github.com/IzuchukwuSamson/lexi/internal/app/users/models"
)

type AdminServiceInterface interface {
	FetchAllUsers() ([]userModel.User, error)
	AdminEmailExists(email string) (bool, error)
	GetUserRoleByID(adminID int) (string, error)
	NewAdminEmail(ad adminModel.Admin) (*adminModel.Admin, error)
}
