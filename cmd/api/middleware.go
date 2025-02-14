package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userCtx contextKey = "user"

func (app *application) GetUserFromCtx(r *http.Request) store.Users {
	user := r.Context().Value(userCtx).(store.Users)
	return user
}

func (app *application) ContextMiddlware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), QueryTimeoutDuration)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("missing authorization header"))
				return
			}
			// parse
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("malformed authorization header"))
				return
			}

			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, err)
			}

			// check creds
			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("missing authorization header"))
				return
			}
			// parse
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("malformed authorization header"))
				return
			}

			token := parts[1]
			jwtToken, err := app.authenticator.ValidateToken(token)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			claims, _ := jwtToken.Claims.(jwt.MapClaims)

			userID, err := uuid.Parse(claims["sub"].(string))
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx := r.Context()

			user, err := app.store.GetUserByUserId(ctx, userID)
			if err != nil {
				app.unauthorizedErrorResponse(w, r, err)
				return
			}

			ctx = context.WithValue(ctx, userCtx, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
