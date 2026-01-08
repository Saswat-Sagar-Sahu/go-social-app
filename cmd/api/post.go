package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type createPostRequest struct {
	Title   string   `json:"title" validate:"required"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	var payLoad createPostRequest
	if err := readJson(r, &payLoad); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := Validate.Struct(payLoad); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payLoad.Title,
		Content: payLoad.Content,
		Tags:    payLoad.Tags,
		// UserID will be NULL for now until auth is implemented
	}

	if err := app.Store.Posts.Create(ctx, post); err != nil {
		log.Println(err)
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
	post, err := app.Store.Posts.GetByID(ctx, postID)
	if err != nil {
		log.Println(err)
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}
	if err := writeJson(w, http.StatusOK, post); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
