package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type createPostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
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
		Tags:    payLoad.Tags,
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

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		writeJsonError(w, http.StatusBadRequest, "invalid post ID")
		return
	}
	post, err := app.Store.Posts.GetByID(ctx, postID)
	if err != nil {
		log.Println(err)
		writeJsonError(w, http.StatusNotFound, "No record found"+postIDStr)
		return
	}
	if err := writeJson(w, http.StatusOK, post); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
