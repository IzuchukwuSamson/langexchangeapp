package handlers

import (
	"log"
	"net/http"

	"github.com/IzuchukwuSamson/lexi/internal/app/admin/services"
)

type AdminHandlers struct {
	log     *log.Logger
	service services.AdminServiceInterface
}

func NewAdminHandlers(l *log.Logger, s services.AdminServiceInterface) *AdminHandlers {
	return &AdminHandlers{
		log:     l,
		service: s,
	}
}

func (a *AdminHandlers) RegisterAdminEmail(rw http.ResponseWriter, r *http.Request) {
	a.log.Println("Register")
}

func (a *AdminHandlers) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	a.log.Println("GET ALL USERS")
}
