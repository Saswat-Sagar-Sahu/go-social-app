package main

import (
	"net/http"

	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

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

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseInt64Param(r, chi.URLParam(r, "userId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload struct {
		FollowerID int64 `json:"follower_id"`
	}
	if err := readJson(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.Store.Followers.AddFollower(ctx, userID, payload.FollowerID); err != nil {
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

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := parseInt64Param(r, chi.URLParam(r, "userId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload struct {
		FollowerID int64 `json:"follower_id"`
	}
	if err := readJson(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := app.Store.Followers.RemoveFollower(ctx, userID, payload.FollowerID); err != nil {
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
