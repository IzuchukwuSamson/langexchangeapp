package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IzuchukwuSamson/lexi/initializer"
	"github.com/IzuchukwuSamson/lexi/internal/db"
	"github.com/IzuchukwuSamson/lexi/router"
	gHandlers "github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file only if not in production
	if os.Getenv("ENVIRONMENT") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file:", err)
		}
	}

	port, host := getHostAndPort()
	logger := log.New(os.Stdout, fmt.Sprintf("%s:", os.Getenv("APP_NAME")), log.LstdFlags)

	dbsql, err := db.SQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbConn := db.NewDB(nil, dbsql)
	if err != nil {
		logger.Fatal(err)
	}

	services := initializer.Services(dbConn)
	handlers := initializer.Handlers(services, logger)

	apiRoutes := router.NewRouter(handlers)

	// CORS middleware
	corsHandler := gHandlers.CORS(
		gHandlers.AllowedOrigins([]string{"http://localhost:3000"}),         // Allow all origins
		gHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),  // Allow specific methods
		gHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow specific headers
		gHandlers.AllowCredentials(),                                        // Allow credentials (cookies, headers, etc.)
	)(apiRoutes)

	logHandler := gHandlers.CombinedLoggingHandler(os.Stdout, corsHandler)

	s := http.Server{
		Addr:         host + ":" + port,
		Handler:      logHandler,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
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
		log.Fatalf("error shutting down server: %v", err)
	}
	log.Print("server shut down gracefully")
}

func getHostAndPort() (string, string) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port if not set
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0" // default host if not set
	}
	return port, host
}
