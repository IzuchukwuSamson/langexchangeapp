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
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	port, host := getHostAndPort()
	logger := log.New(os.Stdout, fmt.Sprintf("%s:", os.Getenv("APP_NAME")), log.LstdFlags)

	dbsql, err := db.SQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	dbConn := db.NewDB(nil, dbsql)
	// redis := db.Redis()
	if err != nil {
		logger.Fatal(err)
	}

	services := initializer.Services(dbConn)
	handlers := initializer.Handlers(services, logger)

	apiRoutes := router.NewRouter(handlers)

	logHandler := gHandlers.CombinedLoggingHandler(os.Stdout, apiRoutes)

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
		log.Fatalf("error shuting down server: %v", err)
	}
	log.Print("server shut down gracefully")
}

func getHostAndPort() (string, string) {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	return port, host
}
