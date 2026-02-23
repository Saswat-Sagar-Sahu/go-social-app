package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type createPostRequest struct {
	Title    string   `json:"title" validate:"required"`
	Content  string   `json:"content" validate:"required,max=1000"`
	ImageURL string   `json:"image_url" validate:"omitempty,url,max=2000"`
	Tags     []string `json:"tags"`
}

type updatePostRequest struct {
	Title    string   `json:"title" validate:"required"`
	Content  string   `json:"content" validate:"required,max=1000"`
	ImageURL string   `json:"image_url" validate:"omitempty,url,max=2000"`
	Tags     []string `json:"tags"`
}

// createPostsHandler handles the creation of a new post
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content, and tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		createPostRequest	true	"Post data"
//	@Success		201		{object}	store.Post
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/posts/ [post]
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
	}
	if payLoad.ImageURL != "" {
		v := payLoad.ImageURL
		post.ImageURL = &v
	}

	uid, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}
	post.UserID = &uid

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

// getPostsHandler retrieves a post by ID
//
//	@Summary		Retrieve a post by ID
//	@Description	Get a post by its unique ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int64	true	"Post ID"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/posts/{postId} [get]
func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
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

// getMyPostsHandler retrieves posts created by the authenticated user.
//
//	@Summary		Retrieve authenticated user's posts
//	@Description	Get all posts created by the authenticated user
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		store.Post
//	@Failure		401	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/posts/me [get]
func (app *application) getMyPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	posts, err := app.Store.Posts.GetByUserID(ctx, uid)
	if err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, posts); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

func (app *application) isAdmin(ctx *http.Request, userID int64) (bool, error) {
	roles, err := app.Store.Users.GetRoles(ctx.Context(), userID)
	if err != nil {
		return false, err
	}
	for _, rr := range roles {
		if rr == "admin" {
			return true, nil
		}
	}
	return false, nil
}

// updatePostHandler updates a post by ID.
//
//	@Summary		Update a post
//	@Description	Update a post by ID (owner or admin)
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int64				true	"Post ID"
//	@Param			post	body		updatePostRequest	true	"Updated post data"
//	@Success		200		{object}	store.Post
//	@Failure		400		{object}	errorResponse
//	@Failure		401		{object}	errorResponse
//	@Failure		403		{object}	errorResponse
//	@Failure		404		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/posts/{postId} [put]
func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var payload updatePostRequest
	if err := readJson(r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, err := app.Store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	uid, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	isOwner := post.UserID != nil && *post.UserID == uid
	if !isOwner {
		admin, err := app.isAdmin(r, uid)
		if err != nil {
			app.statusInternalServerError(w, r, err)
			return
		}
		if !admin {
			app.forbiddenResponse(w, r, ErrInvalidCredentials)
			return
		}
	}

	post.Title = payload.Title
	post.Content = payload.Content
	if payload.ImageURL == "" {
		post.ImageURL = nil
	} else {
		v := payload.ImageURL
		post.ImageURL = &v
	}
	post.Tags = payload.Tags
	if err := app.Store.Posts.Update(ctx, post); err != nil {
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

// deletePostHandler deletes a post by ID.
//
//	@Summary		Delete a post
//	@Description	Delete a post by ID (owner or admin)
//	@Tags			posts
//	@Param			postId	path	int64	true	"Post ID"
//	@Success		204
//	@Failure		400	{object}	errorResponse
//	@Failure		401	{object}	errorResponse
//	@Failure		403	{object}	errorResponse
//	@Failure		404	{object}	errorResponse
//	@Failure		500	{object}	errorResponse
//	@Security		BearerAuth
//	@Router			/posts/{postId} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID, err := strconv.ParseInt(chi.URLParam(r, "postId"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, err := app.Store.Posts.GetByID(ctx, postID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	uid, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	isOwner := post.UserID != nil && *post.UserID == uid
	if !isOwner {
		admin, err := app.isAdmin(r, uid)
		if err != nil {
			app.statusInternalServerError(w, r, err)
			return
		}
		if !admin {
			app.forbiddenResponse(w, r, ErrInvalidCredentials)
			return
		}
	}

	if err := app.Store.Posts.DeleteByID(ctx, postID); err != nil {
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
