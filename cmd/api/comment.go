package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Saswat-Sagar-Sahu/Social/internal/auth"
	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type createCommentRequest struct {
	Content string `json:"content" validate:"required,max=500"`
	PostID  int64  `json:"post_id" validate:"required"`
}

// createCommentsHandler handles creating a comment
// @Summary Create a new comment
// @Description Create a new comment on a post (authenticated)
// @Tags comments
// @Accept json
// @Produce json
// @Param post body createCommentRequest true "Comment data"
// @Success 201 {object} store.Comment
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /comments/ [post]
func (app *application) createCommentsHandler(w http.ResponseWriter, r *http.Request) {
	var payLoad createCommentRequest
	if err := readJson(r, &payLoad); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := Validate.Struct(payLoad); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// set authenticated user as comment owner
	var uid int64
	if u, ok := auth.UserIDFromContext(ctx); ok {
		uid = u
	} else {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	comment := &store.Comment{
		Content: payLoad.Content,
		PostID:  payLoad.PostID,
		UserID:  uid,
	}

	// validate referenced Post exists
	if _, err := app.Store.Posts.GetByID(ctx, payLoad.PostID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	// validate authenticated user exists
	if _, err := app.Store.Users.GetByID(ctx, uid); err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	if err := app.Store.Comments.Create(ctx, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusCreated, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// getCommentsByPostIdHandler retrieves comments for a specific post
//
//	@Summary		Retrieve comments by post ID
//	@Description	Get comments for a specific post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int64	true	"Post ID"
//	@Success		200		{array}		store.Comment
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/comments/post/{postId} [get]
func (app *application) getCommentsByPostIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	postID, err := parseInt64Param(r, chi.URLParam(r, "postId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	storeErr, comments := app.Store.Comments.GetByPostId(ctx, postID)
	if storeErr != nil {
		switch storeErr {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, storeErr)
		default:
			app.statusInternalServerError(w, r, storeErr)
		}
		return
	}

	if err := writeJson(w, http.StatusOK, comments); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

func parseInt64Param(_ *http.Request, s string) (int64, error) {
	postId, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return postId, nil
}

// getCommentsByUserIdHandler retrieves comments made by a specific user
//
//	@Summary		Retrieve comments by user ID
//	@Description	Get comments made by a specific user
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int64	true	"User ID"
//	@Success		200		{array}		store.Comment
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/comments/user/{userId} [get]
func (app *application) getCommentsByUserIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commentID, err := parseInt64Param(r, chi.URLParam(r, "userId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err, comment := app.Store.Comments.GetByUserId(ctx, commentID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}
	if err := writeJson(w, http.StatusOK, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}

// deleteCommentByIdHandler deletes a comment by its ID
//
//	@Summary		Delete a comment by ID
//	@Description	Delete a comment by its unique ID
//	@Tags			comments
//	@Param			commentId	path		int64	true	"Comment ID"
//	@Success		204
//	@Failure		400		{object}	errorResponse
//	@Failure		500		{object}	errorResponse
//	@Router			/comments/{commentId} [delete]
func (app *application) deleteCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commentID, err := parseInt64Param(r, chi.URLParam(r, "commentId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// authorization: only owner or admin can delete
	_, comment := app.Store.Comments.GetByCommentId(ctx, commentID)
	if comment == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	uid, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}

	allowed := false
	if comment.UserID == uid {
		allowed = true
	} else {
		// check admin role
		roles, err := app.Store.Users.GetRoles(ctx, uid)
		if err == nil {
			for _, rr := range roles {
				if rr == "admin" {
					allowed = true
					break
				}
			}
		}
	}
	if !allowed {
		app.forbiddenResponse(w, r, ErrInvalidCredentials)
		return
	}

	if err := app.Store.Comments.DeleteByID(ctx, commentID); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updateCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	var payLoad createCommentRequest
	commentID, err := parseInt64Param(r, chi.URLParam(r, "commentId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := readJson(r, &payLoad); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()
	if err := Validate.Struct(payLoad); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var comment *store.Comment
	err, comment = app.Store.Comments.GetByCommentId(ctx, commentID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	// authorization: only owner or admin can update
	uid, ok := auth.UserIDFromContext(ctx)
	if !ok {
		app.unauthorizedResponse(w, r, ErrInvalidCredentials)
		return
	}
	if comment.UserID != uid {
		roles, err := app.Store.Users.GetRoles(ctx, uid)
		if err != nil {
			app.statusInternalServerError(w, r, err)
			return
		}
		isAdmin := false
		for _, rr := range roles {
			if rr == "admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			app.forbiddenResponse(w, r, ErrInvalidCredentials)
			return
		}
	}

	// apply incoming fields to the loaded comment (do not change owner)
	comment.Content = payLoad.Content
	comment.PostID = payLoad.PostID

	// validate referenced Post exists
	if _, err := app.Store.Posts.GetByID(ctx, payLoad.PostID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, fmt.Errorf("post_id %d is invalid", payLoad.PostID))
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	// owner remains unchanged; no user validation required

	if err := app.Store.Comments.UpdateComment(ctx, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
