package routes

import (
	"lexibuddy/config/middleware"
	"lexibuddy/handlers/users"

	"github.com/gorilla/mux"
)

func NewRouter(userHandlers *users.UserHandlers) *mux.Router {
	router := mux.NewRouter()

	// Define the routes
	router.HandleFunc("/signup", userHandlers.Signup).Methods("POST")
	router.HandleFunc("/login", userHandlers.PasswordLogin).Methods("POST")
	router.HandleFunc("/oauth", userHandlers.SocialLogin).Methods("POST")
	router.HandleFunc("/verify-email", userHandlers.VerifyEmail).Methods("POST")
	router.HandleFunc("/forgot-password", userHandlers.ForgotPassword).Methods("POST")
	router.HandleFunc("/reset-password", userHandlers.ResetPassword).Methods("POST")

	authRoutes := router.NewRoute().Subrouter()
	authRoutes.Use(mux.MiddlewareFunc(middleware.AuthCheck))

	authRoutes.HandleFunc("/getallusers", userHandlers.GetAllUsers).Methods("GET")

	return router
}
