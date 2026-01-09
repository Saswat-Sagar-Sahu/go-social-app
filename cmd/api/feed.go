package main

import "net/http"

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