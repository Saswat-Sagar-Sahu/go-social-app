package main

import (
	"net/http"

	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
)

// forbiddenResponse returns a 403 JSON error
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	writeJsonError(w, http.StatusForbidden, err.Error())
}

// RequireRole returns middleware that allows only users with the given role.
func (app *application) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid, ok := auth.UserIDFromContext(r.Context())
			if !ok {
				app.unauthorizedResponse(w, r, ErrInvalidCredentials)
				return
			}

			roles, err := app.Store.Users.GetRoles(r.Context(), uid)
			if err != nil {
				app.statusInternalServerError(w, r, err)
				return
			}
			for _, rr := range roles {
				if rr == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			app.forbiddenResponse(w, r, ErrInvalidCredentials)
		})
	}
}

// RequireOwnerOrRole returns middleware that allows the request if the authenticated user
// is the owner (determined by getOwnerID) or has the specified role.
func (app *application) RequireOwnerOrRole(getOwnerID func(r *http.Request) (int64, error), role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid, ok := auth.UserIDFromContext(r.Context())
			if !ok {
				app.unauthorizedResponse(w, r, ErrInvalidCredentials)
				return
			}

			ownerID, err := getOwnerID(r)
			if err == nil && ownerID == uid {
				next.ServeHTTP(w, r)
				return
			}

			// check role
			roles, err := app.Store.Users.GetRoles(r.Context(), uid)
			if err != nil {
				app.statusInternalServerError(w, r, err)
				return
			}
			for _, rr := range roles {
				if rr == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			app.forbiddenResponse(w, r, ErrInvalidCredentials)
		})
	}
}
