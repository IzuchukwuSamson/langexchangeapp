package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/IzuchukwuSamson/lexi/internal/app/users/models"
	"github.com/IzuchukwuSamson/lexi/internal/app/users/services"
	"github.com/IzuchukwuSamson/lexi/utils"
	"golang.org/x/oauth2"
)

type UserHandlers struct {
	log      *log.Logger
	services services.UserServiceInterface
}

func NewUserHandlers(l *log.Logger, s services.UserServiceInterface) *UserHandlers {
	return &UserHandlers{
		log:      l,
		services: s,
	}
}

func (u UserHandlers) Signup(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := utils.FromJSON(r.Body, &user)
	if err != nil {
		u.log.Printf("error decoding json request: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "there was an error with your request"}, http.StatusBadRequest)
		return
	}

	user.FirstName = utils.SanitizeInput(user.FirstName)
	user.LastName = utils.SanitizeInput(user.LastName)
	user.Email = strings.ToLower(utils.SanitizeInput(user.Email))

	if user.FirstName == "" || user.LastName == "" || user.Email == "" || user.Password == "" || user.PhoneNumber == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "complete all fields"}, http.StatusBadRequest)
		return
	}

	if !utils.IsValidEmail(user.Email) {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid email format"}, http.StatusBadRequest)
		return
	}

	if !utils.IsValidPassword(user.Password) {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "password must contain uppercase, lowercase, and numbers"}, http.StatusBadRequest)
		return
	}

	existingUser, err := u.services.GetUserByEmail(user.Email)
	if err == nil && existingUser.Email != "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "user exists"}, http.StatusBadRequest)
		return
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error hashing password"}, http.StatusBadRequest)
		return
	}
	user.Role = utils.UserRole
	user.LastActive = time.Now()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Attempt to create the user
	createdUser, verificationCode, err := u.services.CreateUser(user)
	if err != nil {
		u.log.Printf("error creating user: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "there was an error with your request"}, http.StatusBadRequest)
		return
	}

	utilsUser := utils.User{
		FirstName: createdUser.FirstName,
		Email:     createdUser.Email,
	}

	// Send verification email with PIN
	errChan := utils.SendEmail(utils.VerifyEmail, utilsUser, verificationCode)
	if err := <-errChan; err != nil {
		u.log.Printf("error sending verification email: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not send verification email"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw,
		utils.ResponseMsg{
			Message: "Account created successfully. Please check your email for the verification PIN.",
			Data: map[string]string{
				"id": createdUser.ID,
			},
		},
		http.StatusOK,
	)
}

func (u UserHandlers) VerifyEmail(rw http.ResponseWriter, r *http.Request) {
	// var user models.User
	var verificationRequest struct {
		Code string `json:"code"`
	}

	// Decode the request body into the verification request struct
	err := utils.FromJSON(r.Body, &verificationRequest)
	if err != nil {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid request format"}, http.StatusBadRequest)
		return
	}

	// Call the VerifyEmail service method with the verification code
	err = u.services.VerifyEmail(verificationRequest.Code)
	if err != nil {
		errorMessage := "invalid verification code"
		if err.Error() == "invalid or expired verification code" {
			errorMessage = err.Error()
		}
		utils.ReturnJSON(rw, utils.ErrMessage{Error: errorMessage}, http.StatusBadRequest)
		return
	}

	utils.ReturnJSON(rw,
		utils.MailResponse{
			Message: "Email verified successfully",
		},
		http.StatusOK,
	)

}

func (u UserHandlers) PasswordLogin(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := utils.FromJSON(r.Body, &user); err != nil {
		u.log.Printf("error reading login request: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid request"}, http.StatusUnprocessableEntity)
		return
	}

	if user.Email == "" || user.Password == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "complete all fields"}, http.StatusUnprocessableEntity)
		return
	}

	userDb, err := u.services.GetUserByEmail(user.Email)
	if err != nil || userDb == nil {
		u.log.Printf("get user error: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "details not found"}, http.StatusUnauthorized)
		return
	}

	if userDb.ID == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "details not found"}, http.StatusUnauthorized)
		return
	}

	if userDb.EmailVerified == 0 { // 0 means inactive
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "verify your account before logging in"}, http.StatusUnauthorized)
		return
	}

	// Verify the password
	err = utils.ComparePassword(userDb.Password, user.Password)
	if err != nil {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid credentials"}, http.StatusUnauthorized)
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := utils.GenerateToken(fmt.Sprint(userDb.ID))
	if err != nil {
		u.log.Printf("error generating tokens: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not generate tokens"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw, utils.ResponseMsg{
		Message: "Logged in successfully",
		Data: map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		}},
		http.StatusOK,
	)
}

func (u UserHandlers) RefreshToken(rw http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := utils.FromJSON(r.Body, &request); err != nil {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid request"}, http.StatusUnprocessableEntity)
		return
	}

	newAccessToken, err := utils.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid refresh token"}, http.StatusUnauthorized)
		return
	}

	utils.ReturnJSON(rw, utils.ResponseMsg{
		Message: "Token refreshed successfully",
		Data: map[string]string{
			"access_token": newAccessToken,
		}},
		http.StatusOK,
	)
}

func (u UserHandlers) ForgotPassword(rw http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	if err := utils.FromJSON(r.Body, &request); err != nil {
		u.log.Printf("error reading forgot password request: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid request"}, http.StatusUnprocessableEntity)
		return
	}

	if request.Email == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "email is required"}, http.StatusUnprocessableEntity)
		return
	}

	user, err := u.services.GetUserByEmail(request.Email)
	if err != nil || user == nil {
		u.log.Printf("get user error: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "email not found"}, http.StatusNotFound)
		return
	}

	resetToken, err := u.services.GeneratePasswordResetToken(*user)
	if err != nil {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not process password reset request"}, http.StatusInternalServerError)
		return
	}

	utilsUser := utils.User{
		FirstName: user.FirstName,
		Email:     user.Email,
	}

	errChan := utils.SendEmail(utils.ResetPassword, utilsUser, resetToken)
	if err := <-errChan; err != nil {
		u.log.Printf("error sending reset pin: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not send reset pin"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw, utils.MailResponse{
		Message: "Password reset email sent successfully",
	}, http.StatusOK)
}

func (u UserHandlers) ResetPassword(rw http.ResponseWriter, r *http.Request) {
	var request struct {
		Email           string `json:"email"`
		ResetCode       string `json:"reset_code"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	err := utils.FromJSON(r.Body, &request)
	if err != nil {
		u.log.Printf("error decoding json request: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid request format"}, http.StatusBadRequest)
		return
	}

	if request.Email == "" || request.ResetCode == "" || request.NewPassword == "" || request.ConfirmPassword == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "complete all fields"}, http.StatusBadRequest)
		return
	}

	if request.NewPassword != request.ConfirmPassword {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "passwords do not match"}, http.StatusBadRequest)
		return
	}

	user, err := u.services.GetUserByEmail(request.Email)
	if err != nil {
		u.log.Printf("error finding user by email: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "user not found"}, http.StatusNotFound)
		return
	}

	userIDInt, err := strconv.ParseInt(user.ID, 10, 64)
	if err != nil {
		u.log.Printf("error parsing user ID: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid user ID"}, http.StatusInternalServerError)
		return
	}

	passwordReset, err := u.services.GetPasswordResetByCode(request.ResetCode)
	if err != nil || passwordReset == nil || passwordReset.ExpiresAt.Before(time.Now()) || passwordReset.UserID != userIDInt {
		u.log.Printf("invalid or expired reset code: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid or expired reset code"}, http.StatusUnauthorized)
		return
	}

	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		u.log.Printf("error hashing password: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error processing new password"}, http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	err = u.services.UpdateUserPassword(user)
	if err != nil {
		u.log.Printf("error updating user password: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not update password"}, http.StatusInternalServerError)
		return
	}

	// Delete the password reset record from the database
	err = u.services.DeletePasswordReset(passwordReset.ID)
	if err != nil {
		u.log.Printf("error deleting password reset record: %v\n", err)
	}

	utils.ReturnJSON(rw,
		utils.MailResponse{
			Message: "Password reset successfully",
		},
		http.StatusOK,
	)
}

func (u UserHandlers) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	users, err := u.services.FetchAllUsers()
	if err != nil {
		u.log.Printf("error decoding json request: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error getting users"}, http.StatusInternalServerError)
		return
	}

	// dto := make([]models.UserDTO, len(users))
	// for i, user := range users {
	// 	dto[i] = utils.ToUserDTO(user)
	// }

	utils.ReturnJSON(rw, users, http.StatusOK)
}

func (u UserHandlers) Logout(rw http.ResponseWriter, r *http.Request) {
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "no token provided"}, http.StatusUnauthorized)
		return
	}

	tokenString := strings.Split(authHeader, "Bearer ")[1]
	if tokenString == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "no token provided"}, http.StatusUnauthorized)
		return
	}

	// Invalidate the token (depends on your implementation, e.g., adding to a blacklist)
	err := u.services.InvalidateToken(tokenString)
	if err != nil {
		u.log.Printf("error invalidating token: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not logout"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw, map[string]string{"message": "logged out successfully"}, http.StatusOK)
}

func (u UserHandlers) SocialLogin(rw http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	code := r.URL.Query().Get("code")
	var config *oauth2.Config

	switch provider {
	case "google":
		config = utils.GoogleConfig
	case "facebook":
		config = utils.FacebookConfig
	// case "apple":
	// 	config = appleConfig
	default:
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "unknown provider"}, http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		u.log.Printf("token exchange error: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "token exchange failed"}, http.StatusUnauthorized)
		return
	}

	userInfo, err := fetchUserInfo(provider, token)
	if err != nil {
		u.log.Printf("fetch user info error: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not fetch user info"}, http.StatusUnauthorized)
		return
	}

	user, err := u.services.FindOrCreateUser(userInfo)
	if err != nil {
		u.log.Printf("find or create user error: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not log in user"}, http.StatusInternalServerError)
		return
	}

	tokenString, err := utils.GenerateJWT(user)
	if err != nil {
		u.log.Printf("error generating token: %v\n", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "could not generate token"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw, map[string]string{"token": tokenString}, http.StatusOK)
}

func fetchUserInfo(provider string, token *oauth2.Token) (map[string]interface{}, error) {
	var endpoint string

	switch provider {
	case "google":
		endpoint = "https://www.googleapis.com/oauth2/v2/userinfo"
	case "facebook":
		endpoint = "https://graph.facebook.com/me?fields=id,name,email"
	case "apple":
		endpoint = "https://appleid.apple.com/auth/keys"
	default:
		return nil, errors.New("unknown provider")
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
