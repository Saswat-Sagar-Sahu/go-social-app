package main

import (
	"log"
	"net/http"

	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
)

type createPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (app *application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	var payLoad createPostRequest
	if err := readJson(r, &payLoad); err != nil {
		writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()

	post := &store.Post{
		Title:   payLoad.Title,
		Content: payLoad.Content,
		// UserID will be NULL for now until auth is implemented
	}

	if err := app.Store.Posts.Create(ctx, post); err != nil {
		log.Println(err)
		writeJsonError(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	if err := writeJson(w, http.StatusCreated, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
