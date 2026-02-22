package main

import (
	"net/http"

	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
)

// getUserFeedHandler retrieves the feed for a specific user
//
//	@Summary		Retrieve user feed
//	@Description	Get the feed for a specific user
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Post
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	feed, err := app.Store.Posts.GetUserFeed(ctx, userID)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, feed); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
