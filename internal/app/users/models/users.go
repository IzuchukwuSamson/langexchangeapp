package models

import (
	"time"
)

type User struct {
	ID                      int64                  `json:"id"`
	FirstName               string                 `json:"firstname"`
	LastName                string                 `json:"lastname"`
	Username                string                 `json:"username"`
	Email                   string                 `json:"email"`
	EmailVerified           int                    `json:"email_verified"`
	Password                string                 `json:"password"`
	PasswordConfirm         string                 `json:"password_confirm"`
	PhoneNumber             string                 `json:"phonenumber"`
	Role                    string                 `json:"role"`
	LastActive              time.Time              `json:"last_active"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
	DeletedAt               time.Time              `json:"deleted_at"`
	LanguagesSpoken         []string               `json:"languages_spoken"`
	LanguagesLearning       []string               `json:"languages_learning"`
	Bio                     string                 `json:"bio"`
	ProfilePicture          string                 `json:"profile_picture"`
	Location                string                 `json:"location"`
	LearningGoals           []string               `json:"learning_goals"`
	PreferredLanguage       string                 `json:"preferred_language"`
	TimeZone                string                 `json:"time_zone"`
	AccountType             string                 `json:"account_type"`
	SocialLinks             map[string]string      `json:"social_links"`
	Connections             []string               `json:"connections"`
	PreferredLearningStyle  string                 `json:"preferred_learning_style"`
	Availability            string                 `json:"availability"`
	Achievements            []string               `json:"achievements"`
	Progress                map[string]interface{} `json:"progress"`
	Ratings                 map[string]float64     `json:"ratings"`
	Feedback                []string               `json:"feedback"`
	SubscriptionStatus      string                 `json:"subscription_status"`
	PaymentMethod           string                 `json:"payment_method"`
	NotificationPreferences map[string]bool        `json:"notification_preferences"`
	TwoFactorEnabled        bool                   `json:"two_factor_enabled"`
	PrivacySettings         map[string]bool        `json:"privacy_settings"`
}

type UserDTO struct {
	ID          int64     `json:"id" bson:"_id,omitempty"`
	FirstName   string    `json:"firstname" bson:"firstname"`
	LastName    string    `json:"lastname" bson:"lastname"`
	Email       string    `json:"email" bson:"email"`
	PhoneNumber string    `json:"phonenumber" bson:"phonenumber"`
	Role        string    `json:"role" bson:"role"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type EmailVerification struct {
	ID     int64  `json:"id"`
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
