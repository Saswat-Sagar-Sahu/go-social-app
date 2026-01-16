package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
)

func (app *application) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedResponse(w, r, ErrInvalidCredentials)
			return
		}
		parts := strings.Fields(authHeader)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			app.unauthorizedResponse(w, r, ErrInvalidCredentials)
			return
		}
		token := parts[1]
		uid, err := app.Auth.ParseToken(r.Context(), token)
		if err != nil {
			app.unauthorizedResponse(w, r, ErrInvalidCredentials)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextUserKey, uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// methodAuthMiddleware enforces authMiddleware for POST, PUT, DELETE methods
func (app *application) methodAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// whitelist endpoints that should remain public
		path := r.URL.Path
		if path == "/v1/users/register" || path == "/v1/users/login" || path == "/v1/users/activate" || path == "/v1/health" {
			next.ServeHTTP(w, r)
			return
		}
		// allow swagger paths
		if len(path) >= 8 && path[:8] == "/swagger" {
			next.ServeHTTP(w, r)
			return
		}

		// require auth for everything else
		app.authMiddleware(next).ServeHTTP(w, r)
		return
	})
}
