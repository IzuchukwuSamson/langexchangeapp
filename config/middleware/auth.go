package middleware

import (
	"context"
	"errors"
	"fmt"
	"lexibuddy/services"
	"lexibuddy/utils"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

// AuthCheck middleware performs an authentication check on the request
//
// It checks the header for the Bearer Authorization token and validates it with utils.VerifyToken.
// Then it creates a context with the authenticated object
func AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(utils.AuthHeader)
			lenSchema := len(utils.AuthHeaderPrefix)
			if lenSchema >= len(authHeader) {
				utils.ReturnJSON(rw, utils.ErrMessage{Error: "auth header is needed"}, http.StatusUnauthorized)
				return
			}
			token := authHeader[lenSchema:]
			// u.logger.Println("token: ", token)

			claims, err := utils.VerifyToken(token)
			if err != nil {
				fmt.Println("error verifying token: ", err)
				msg := "user not authorized"
				switch {
				case errors.Is(err, jwt.ErrTokenExpired):
					msg = jwt.ErrTokenExpired.Error()
				}
				utils.ReturnJSON(rw, utils.ErrMessage{Error: msg}, http.StatusUnauthorized)
				return
			}
			id, err := claims.GetSubject()
			if err != nil {
				fmt.Println("error getting claims: ", err)
				utils.ReturnJSON(rw, utils.ErrMessage{Error: utils.GenericError}, http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), utils.CtxKey("user"), id)
			// fmt.Println("user id: ", id)
			req := r.WithContext(ctx)
			next.ServeHTTP(rw, req)
		},
	)
}

func JWTMiddleware(next http.Handler, userService *services.UserService) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ReturnJSON(rw, utils.ErrMessage{Error: "no token provided"}, http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]
		if tokenString == "" {
			utils.ReturnJSON(rw, utils.ErrMessage{Error: "no token provided"}, http.StatusUnauthorized)
			return
		}

		// Check if the token is invalidated
		if userService.IsTokenInvalidated(tokenString) {
			utils.ReturnJSON(rw, utils.ErrMessage{Error: "token is invalidated"}, http.StatusUnauthorized)
			return
		}

		// Add token validation logic here (e.g., parse and validate JWT)
		// ...

		next.ServeHTTP(rw, r)
	})
}
