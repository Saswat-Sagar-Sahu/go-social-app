package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type listUsersResponse struct {
	Data       []*store.User `json:"data"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	Total      int           `json:"total"`
	TotalPages int           `json:"total_pages"`
}

// listUsersHandler retrieves paginated users with optional username filtering.
//
//	@Summary		List users
//	@Description	Get paginated users. Optional `name` filters by username.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int		false	"Page number (default 1)"
//	@Param			page_size	query		int		false	"Page size (default 10, max 100)"
//	@Param			name		query		string	false	"Filter by username (contains match)"
//	@Success		200			{object}	listUsersResponse
//	@Failure		400			{object}	errorResponse
//	@Failure		500			{object}	errorResponse
//	@Router			/users [get]
func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()

	page := 1
	pageSize := 10

	if raw := strings.TrimSpace(q.Get("page")); raw != "" {
		p, err := strconv.Atoi(raw)
		if err != nil || p < 1 {
			app.badRequestResponse(w, r, errors.New("page must be a positive integer"))
			return
		}
		page = p
	}
	if raw := strings.TrimSpace(q.Get("page_size")); raw != "" {
		s, err := strconv.Atoi(raw)
		if err != nil || s < 1 || s > 100 {
			app.badRequestResponse(w, r, errors.New("page_size must be between 1 and 100"))
			return
		}
		pageSize = s
	}

	name := strings.TrimSpace(q.Get("name"))
	if name == "" {
		// allow alias query param for clients that prefer explicit field naming
		name = strings.TrimSpace(q.Get("username"))
	}

	users, total, err := app.Store.Users.List(ctx, store.UserFilter{
		Page:     page,
		PageSize: pageSize,
		Name:     name,
	})
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(pageSize)))
	}

	resp := listUsersResponse{
		Data:       users,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
	if err := writeJson(w, http.StatusOK, resp); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// GetUsersHandler retrieves a user by ID
//
//	@Summary		Retrieve a user by ID
//	@Description	Get a user by their unique ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int64	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/users/{userId} [get]
func (app *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseInt64Param(r, chi.URLParam(r, "userId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.Store.Users.GetByID(ctx, userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	if err := writeJson(w, http.StatusOK, user); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// registerUserHandler handles user registration
//
//	@Summary		Register a new user
//	@Description	Register a new user with username, email, and password
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		object{username=string,email=string,password=string}	true	"User registration data"
//	@Success		201		{object}	map[string]interface{}
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/users/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJson(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Username == "" || payload.Email == "" || payload.Password == "" {
		app.badRequestResponse(w, r, ErrInvalidRequest)
		return
	}

	hashed, err := app.Auth.HashPassword(payload.Password)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashed),
	}
	if err := app.Store.Users.Create(ctx, user); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	// generate activation token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
	token := hex.EncodeToString(b)

	expiry := time.Now().Add(time.Duration(app.config.activationTokenExpiryMinutes) * time.Minute)
	t := &store.Token{UserID: user.ID, Token: token, ExpiresAt: expiry}
	if err := app.Store.Tokens.Create(ctx, t); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	// for now, log the token instead of sending email
	log.Printf("Activation token for %s: %s", user.Email, token)

	if err := writeJson(w, http.StatusCreated, map[string]interface{}{"id": user.ID, "email": user.Email}); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// loginHandler handles user login and returns a JWT
//
// @Summary Login
// @Description Authenticate a user and return a JWT
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object{email=string,password=string} true "Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 403 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /users/login [post]
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := readJson(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Email == "" || payload.Password == "" {
		app.badRequestResponse(w, r, ErrInvalidRequest)
		return
	}

	user, err := app.Store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}
	if !user.Activated {
		app.forbiddenResponse(w, r, ErrAccountNotActivated)
		return
	}

	if err := app.Auth.ComparePassword(user.Password, payload.Password); err != nil {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	token, err := app.Auth.GenerateToken(ctx, user.ID)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, map[string]string{"token": token}); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// activateUserHandler handles user account activation
//
//	@Summary		Activate a user account
//	@Description	Activate a user account using the provided token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			activation	body		object{token=string}	true	"Activation token"
//	@Success		202
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/users/activate [post]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload struct {
		Token string `json:"token"`
	}
	if err := readJson(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if payload.Token == "" {
		app.badRequestResponse(w, r, ErrInvalidRequest)
		return
	}

	t, err := app.Store.Tokens.GetByToken(ctx, payload.Token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	if time.Now().After(t.ExpiresAt) {
		app.badRequestResponse(w, r, ErrInvalidRequest)
		return
	}

	if err := app.Store.Users.Activate(ctx, t.UserID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	if err := app.Store.Tokens.DeleteByToken(ctx, payload.Token); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// followUserHandler handles following a user
//
//	@Summary		Follow a user
//	@Description	Follow a user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int64	true	"User ID to follow"
//	@Success		204
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		409		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/users/{userId}/follow [post]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseInt64Param(r, chi.URLParam(r, "userId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	followerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}
	if followerID == userID {
		app.badRequestResponse(w, r, errors.New("users cannot follow themselves"))
		return
	}
	if err := app.Store.Followers.AddFollower(ctx, userID, followerID); err != nil {
		switch err {
		case store.ErrResourceExists:
			app.conflictResponse(w, r, err)
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// unfollowUserHandler handles unfollowing a user
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int64	true	"User ID to unfollow"
//	@Success		204
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/users/{userId}/unfollow [post]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseInt64Param(r, chi.URLParam(r, "userId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	followerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}
	if err := app.Store.Followers.RemoveFollower(ctx, userID, followerID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
