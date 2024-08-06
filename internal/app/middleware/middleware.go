package middleware

import (
	"log"
)

type Middleware struct {
	logger *log.Logger
}

// NewMiddleware initializes a new Middleware handler
func NewMiddleware(log *log.Logger) *Middleware {
	return &Middleware{
		logger: log,
	}
}
