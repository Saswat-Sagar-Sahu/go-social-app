package main

import "net/http"

// getUserFeedHandler retrieves the feed for a specific user
//
//	@Summary		Retrieve user feed
//	@Description	Get the feed for a specific user
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Post
//	@Failure		500	{object}	errorResponse
//	@Router			/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feed, err := app.Store.Posts.GetUserFeed(ctx, int64(2))
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, feed); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
