package initializer

import (
	"log"

	"github.com/IzuchukwuSamson/lexi/internal/app/middleware"
	"github.com/IzuchukwuSamson/lexi/internal/app/users/handlers"
)

type Handler struct {
	User       handlers.UserHandlers
	Middleware middleware.Middleware
}

func Handlers(services *Store, log *log.Logger) *Handler {
	return &Handler{
		// User:         *handlers.NewUser(log, services.User, services.Redis),
		// Middleware:   *middleware.NewMiddleware(log),
		User: *handlers.NewUserHandlers(log, services.User),
	}
}
