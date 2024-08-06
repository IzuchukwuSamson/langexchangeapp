package router

import (
	"net/http"

	"github.com/IzuchukwuSamson/lexi/initializer"
	"github.com/IzuchukwuSamson/lexi/internal/app/users"
	"github.com/gorilla/mux"
)

// NewRouter creates a new global router and adds routes to it
func NewRouter(handlers *initializer.Handler) *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()

	registerHelloHandler(api)

	// register individual routes
	users.RegisterUserRoutes(api, handlers)

	return router
}

func registerHelloHandler(api *mux.Router) *mux.Route {
	return api.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello World\n"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
