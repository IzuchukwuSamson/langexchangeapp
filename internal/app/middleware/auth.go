package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/IzuchukwuSamson/lexi/utils"
	"github.com/golang-jwt/jwt/v5"
)

// AuthCheck middleware performs an authentication check on the request
//
// It checks the header for the Bearer Authorization token and validates it with utils.VerifyToken.
// Then it creates a context with the authenticated object
// func (m *Middleware) AuthCheck(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(rw http.ResponseWriter, r *http.Request) {
// 			authHeader := r.Header.Get(utils.AuthHeader)
// 			lenSchema := len(utils.AuthHeaderPrefix)
// 			if lenSchema >= len(authHeader) {
// 				utils.ReturnJSON(rw, utils.ErrMessage{Error: "auth header is needed"}, http.StatusUnauthorized)
// 				return
// 			}
// 			token := authHeader[lenSchema:]
// 			// u.logger.Println("token: ", token)

// 			claims, err := utils.VerifyToken(token)
// 			if err != nil {
// 				m.logger.Println("error verifying token: ", err)
// 				msg := "user not authorized"
// 				switch {
// 				case errors.Is(err, jwt.ErrTokenExpired):
// 					msg = jwt.ErrTokenExpired.Error()
// 				}
// 				utils.ReturnJSON(rw, utils.ErrMessage{Error: msg}, http.StatusUnauthorized)
// 				return
// 			}
// 			id, err := claims.GetSubject()
// 			if err != nil {
// 				m.logger.Println("error getting claims: ", err)
// 				utils.ReturnJSON(rw, utils.ErrMessage{Error: utils.GenericError}, http.StatusInternalServerError)
// 				return
// 			}
// 			ctx := context.WithValue(r.Context(), utils.CtxKey("user"), id)
// 			// fmt.Println("user id: ", id)
// 			req := r.WithContext(ctx)
// 			next.ServeHTTP(rw, req)
// 		},
// 	)
// }

func (m *Middleware) AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		// m.logger.Println("Authorization header: ", authHeader)
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.ReturnJSON(rw, utils.ErrMessage{Error: "auth header is needed"}, http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.VerifyToken(token)
		if err != nil {
			m.logger.Println("error verifying token: ", err)
			msg := "user not authorized"
			if errors.Is(err, jwt.ErrTokenExpired) {
				msg = jwt.ErrTokenExpired.Error()
			}
			utils.ReturnJSON(rw, utils.ErrMessage{Error: msg}, http.StatusUnauthorized)
			return
		}
		id, err := claims.GetSubject()
		if err != nil {
			m.logger.Println("error getting claims: ", err)
			utils.ReturnJSON(rw, utils.ErrMessage{Error: utils.GenericError}, http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), utils.CtxKey("user"), id)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
