package handlers

import (
	"log"
	"net/http"

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

func (a *AdminHandlers) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	a.log.Println("GET ALL USERS")
}
