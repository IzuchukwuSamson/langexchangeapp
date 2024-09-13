package users

import (
	"github.com/IzuchukwuSamson/lexi/initializer"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers all the routes for accessing the User handler operations
func RegisterUserRoutes(base *mux.Router, userHandlers *initializer.Handler) {
	userRoutes := base.PathPrefix("/users").Subrouter()

	userRoutes.HandleFunc("/signup", userHandlers.User.Signup).Methods("POST")
	userRoutes.HandleFunc("/login", userHandlers.User.PasswordLogin).Methods("POST")
	userRoutes.HandleFunc("/oauth", userHandlers.User.SocialLogin).Methods("POST")
	userRoutes.HandleFunc("/verify-email", userHandlers.User.VerifyEmail).Methods("POST")
	userRoutes.HandleFunc("/forgot-password", userHandlers.User.ForgotPassword).Methods("POST")
	userRoutes.HandleFunc("/reset-password", userHandlers.User.ResetPassword).Methods("POST")
	userRoutes.HandleFunc("/getallusers", userHandlers.User.GetAllUsers).Methods("GET")
	userRoutes.HandleFunc("/getuserbyid", userHandlers.User.GetUserById).Methods("GET")

	// Auth protected routes
	userAuthRoutes := userRoutes.NewRoute().Subrouter()
	userAuthRoutes.Use(mux.MiddlewareFunc(userHandlers.Middleware.AuthCheck))

	// Add the GetCurrentUser route
	userAuthRoutes.HandleFunc("/current", userHandlers.User.GetCurrentUser).Methods("GET")

	// Other auth protected routes
	userAuthRoutes.HandleFunc("/dashboard", userHandlers.User.Dashboard).Methods("GET")
	userAuthRoutes.HandleFunc("/user-test", userHandlers.User.UserTestRoute).Methods("GET")
}
