package models

import (
	"time"
)

type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstname"`
	LastName    string    `json:"lastname"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	PhoneNumber string    `json:"phonenumber"`
	Role        string    `json:"role"`
	IsActive    int64     `json:"is_active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
}

type UserDTO struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	FirstName   string    `json:"firstname" bson:"firstname"`
	LastName    string    `json:"lastname" bson:"lastname"`
	Email       string    `json:"email" bson:"email"`
	PhoneNumber string    `json:"phonenumber" bson:"phonenumber"`
	Role        string    `json:"role" bson:"role"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type EmailVerification struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Code   string `json:"code"`
}

type PasswordReset struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at" `
}

func (User) TableName() string {
	return "users"
}
