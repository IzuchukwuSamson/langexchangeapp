package services

import (
	"lexibuddy/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInterface interface {
	CreateUser(user models.User) (*models.User, string, error)
	FindOrCreateUser(userInfo map[string]interface{}) (*models.User, error)
	FetchAllUsers() ([]models.User, error)
	FetchUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	EditUser(id string, updatedUser models.User) (*models.User, error)
	RemoveUser(id string) (*models.User, error)
	UpdateUser(updatedUser *models.User) (*models.User, error)
	InvalidateToken(token string) error
	VerifyEmail(verificationCode string) error
	GeneratePasswordResetToken(user models.User) (string, error)
	GetPasswordResetByCode(code string) (*models.PasswordReset, error)
	UpdateUserPassword(user *models.User) error
	DeletePasswordReset(id primitive.ObjectID) error
}
