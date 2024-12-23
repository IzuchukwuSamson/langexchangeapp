package admin

import (
	"github.com/IzuchukwuSamson/lexi/initializer"
	"github.com/gorilla/mux"
)

func RegisterAdminRoutes(base *mux.Router, handlers *initializer.Handler) {
	userRoutes := base.PathPrefix("/admin").Subrouter()

	userRoutes.HandleFunc("/register", handlers.Admin.RegisterAdminEmail).Methods("POST")
	userRoutes.HandleFunc("/create", handlers.Admin.UpdateAdmin).Methods("POST")
	userRoutes.HandleFunc("/getallusers", handlers.Admin.GetAllUsers).Methods("GET")

	// auth protected routes
	userAuthRoutes := userRoutes.NewRoute().Subrouter()
	userAuthRoutes.Use(mux.MiddlewareFunc(handlers.Middleware.AuthCheck))
	// userAuthRoutes.HandleFunc("/getallusers", handlers.Admin.GetAllUsers).Methods("GET")
}
