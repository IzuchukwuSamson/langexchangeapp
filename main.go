package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/IzuchukwuSamson/lexi/config/db"
	"github.com/IzuchukwuSamson/lexi/handlers/users"
	"github.com/IzuchukwuSamson/lexi/routes"
	"github.com/IzuchukwuSamson/lexi/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gHandlers "github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port, host := getHostAndPort()
	logger := log.New(os.Stdout, fmt.Sprintf("%s:", os.Getenv("APP_NAME")), log.LstdFlags)

	mongodb, err := db.Mongo()
	if err != nil {
		logger.Fatal(err)
	}

	// initialze db
	dbConn := db.NewDB(mongodb, nil)

	// Initialize the service
	userService := services.NewUserService(dbConn.Mongo, logger)

	// Initialize the handlers
	userHandlers := users.NewUserHandlers(logger, userService)

	// Initialize the router
	// methods := []string{"POST", "GET", "PUT", "DELETE"}
	router := routes.NewRouter(userHandlers)
	logHandler := gHandlers.CombinedLoggingHandler(os.Stdout, router)

	s := http.Server{
		Addr:         host + ":" + port,
		Handler:      logHandler,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
	go func() {
		if err := s.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()
	log.Printf("server started on port %s", port)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	// blocks until this signal is received
	recv := <-sigchan
	log.Printf("signal %v received\n", recv)

	shutdownServer(&s)

}

func shutdownServer(s *http.Server) {
	// shutdown gracefully with 5 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("error shuting down server: %v", err)
	}
	log.Print("server shut down gracefully")
}

func getHostAndPort() (string, string) {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	return port, host
}
