package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/JaskiratAnand/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const userCtx contextKey = "user"
const postCtx contextKey = "post"

func (app *application) GetUserFromCtx(r *http.Request) store.Users {
	user := r.Context().Value(userCtx).(store.Users)
	return user
}

func (app *application) GetPostFromCtx(r *http.Request) store.Posts {
	post := r.Context().Value(postCtx).(store.Posts)
	return post
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

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := app.GetUserFromCtx(r)

		idParam := chi.URLParam(r, "postID")
		postID, err := uuid.Parse(idParam)
		if err != nil {
			err = errors.New("invalid post-id")
			app.badRequestResponse(w, r, err)
			return
		}

		post, err := app.store.GetPostsById(ctx, postID)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), &user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.Users, roleName string) (bool, error) {
	role, err := app.store.GetRoleByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.RoleID >= role.ID, nil
}
