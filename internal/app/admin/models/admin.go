package models

import "time"

type Admin struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	// EmailVerified int       `json:"email_verified"`
	Password     string    `json:"password"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	LastLoggedIn int64     `json:"last_logged_in"`
	Device       string    `json:"device"`
	City         string    `json:"city"`
	Country      string    `json:"country"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}
