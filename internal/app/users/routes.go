package users

import (
	"github.com/IzuchukwuSamson/lexi/initializer"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers all the routes for accessing the User handler operations
func RegisterUserRoutes(base *mux.Router, userHandlers *initializer.Handler) {
	userRoutes := base.PathPrefix("/users").Subrouter()

	userRoutes.HandleFunc("/signup", userHandlers.User.Signup)
	userRoutes.HandleFunc("/login", userHandlers.User.PasswordLogin).Methods("POST")
	userRoutes.HandleFunc("/oauth", userHandlers.User.SocialLogin).Methods("POST")
	userRoutes.HandleFunc("/verify-email", userHandlers.User.VerifyEmail).Methods("POST")
	userRoutes.HandleFunc("/forgot-password", userHandlers.User.ForgotPassword).Methods("POST")
	userRoutes.HandleFunc("/reset-password", userHandlers.User.ResetPassword).Methods("POST")

	// auth protected routes
	userAuthRoutes := userRoutes.NewRoute().Subrouter()
	userAuthRoutes.Use(mux.MiddlewareFunc(userHandlers.Middleware.AuthCheck))
	userAuthRoutes.HandleFunc("/getallusers", userHandlers.User.GetAllUsers).Methods("GET")
}
