package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/IzuchukwuSamson/lexi/internal/app/admin/models"
	"github.com/IzuchukwuSamson/lexi/internal/app/admin/services"
	"github.com/IzuchukwuSamson/lexi/utils"
)

type AdminHandlers struct {
	log      *log.Logger
	services services.AdminServiceInterface
}

func NewAdminHandlers(l *log.Logger, s services.AdminServiceInterface) *AdminHandlers {
	return &AdminHandlers{
		log:      l,
		services: s,
	}
}

func (ah *AdminHandlers) RegisterAdminEmail(rw http.ResponseWriter, r *http.Request) {
	var admin models.Admin
	// Read request body
	if err := utils.FromJSON(r.Body, &admin); err != nil {
		ah.log.Printf("Error reading body: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Invalid request body"}, http.StatusBadRequest)
		return
	}
	// Check if email already exists
	exists, err := ah.services.AdminEmailExists(admin.Email)
	if err != nil {
		ah.log.Printf("Error checking email existence: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Error checking email existence"}, http.StatusInternalServerError)
		return
	}
	if exists {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Email already exists"}, http.StatusBadRequest)
		return
	}
	// Generate verification token and expiry date
	token, err := utils.GenerateRandomNumber()
	if err != nil {
		ah.log.Printf("Error generating token: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Error generating token"}, http.StatusInternalServerError)
		return
	}
	// Prepare email data and send email
	emailData := utils.EmailInfo{
		Email: admin.Email,
	}

	errChan := utils.SendEmail(utils.SendAdminLink, emailData, token)
	if err := <-errChan; err != nil {
		ah.log.Printf("Error sending verification email: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Could not send verification email"}, http.StatusInternalServerError)
		return
	}

	// Create new admin email record
	created, err := ah.services.NewAdminEmail(admin)
	if err != nil {
		ah.log.Printf("Error creating new admin email: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "Could not create new admin email"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw,
		utils.ResponseMsg{
			Message: "Admin created successfully",
			Data: map[string]string{
				"id":    created.ID,
				"email": created.Email,
			},
		},
		http.StatusCreated,
	)
}

func (ah *AdminHandlers) UpdateAdmin(rw http.ResponseWriter, r *http.Request) {
	var admin models.Admin

	// Parse the request body into the admin struct
	if err := utils.FromJSON(r.Body, &admin); err != nil {
		log.Printf("error reading body: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error in request body"}, http.StatusBadRequest)
		return
	}

	// Trim and validate the email and password
	admin.Email = strings.TrimSpace(admin.Email)

	if admin.Email == "" || admin.Password == "" {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "email or password is missing"}, http.StatusBadRequest)
		return
	}

	if !utils.IsValidEmail(admin.Email) {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "invalid email format"}, http.StatusBadRequest)
		return
	}

	if !utils.IsValidPassword(admin.Password) {
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "password must contain uppercase, lowercase, and numbers"}, http.StatusBadRequest)
		return
	}

	// Check if the admin already exists
	retrieved, err := ah.services.GetAdminByEmail(admin.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		// If there's an error that's not "no rows found", return an internal server error
		log.Printf("error retrieving admin by email: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error retrieving admin"}, http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashed, err := utils.HashPassword(admin.Password)
	if err != nil {
		errMsg := "error hashing password"
		log.Printf("%v: %v", errMsg, err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: errMsg}, http.StatusInternalServerError)
		return
	}

	// If the admin exists, update the password
	retrieved.Password = hashed
	if err := ah.services.UpdateAdminPassword(retrieved); err != nil {
		log.Printf("error updating admin password: %v", err)
		utils.ReturnJSON(rw, utils.ErrMessage{Error: "error updating admin password"}, http.StatusInternalServerError)
		return
	}

	utils.ReturnJSON(rw,
		utils.ResponseMsg{
			Message: "Admin password updated successfully",
			Data: map[string]interface{}{
				"Admin": retrieved,
			},
		},
		http.StatusOK,
	)
}

func (a *AdminHandlers) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	a.log.Println("GET ALL USERS")
}
