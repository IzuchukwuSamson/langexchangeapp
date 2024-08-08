package initializer

import (
	"log"

	adminHandlers "github.com/IzuchukwuSamson/lexi/internal/app/admin/handlers"
	"github.com/IzuchukwuSamson/lexi/internal/app/middleware"
	userHandlers "github.com/IzuchukwuSamson/lexi/internal/app/users/handlers"
)

type Handler struct {
	User       userHandlers.UserHandlers
	Admin      adminHandlers.AdminHandlers
	Middleware middleware.Middleware
}

func Handlers(services *Store, log *log.Logger) *Handler {
	return &Handler{
		// User:         *handlers.NewUser(log, services.User, services.Redis),
		// Middleware:   *middleware.NewMiddleware(log),
		User:  *userHandlers.NewUserHandlers(log, services.User),
		Admin: *adminHandlers.NewAdminHandlers(log, services.Admin),
	}
}
