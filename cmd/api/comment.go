package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Saswat-Sagar-Sahu/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type createCommentRequest struct {
	Content string `json:"content" validate:"required,max=500"`
	PostID  int64  `json:"post_id" validate:"required"`
	UserID  int64  `json:"user_id" validate:"required"`
}

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

	comment := &store.Comment{
		Content: payLoad.Content,
		PostID:  payLoad.PostID,
		UserID:  payLoad.UserID,
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

	// validate referenced User exists
	if _, err := app.Store.Users.GetByID(ctx, payLoad.UserID); err != nil {
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

func (app *application) deleteCommentByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	commentID, err := parseInt64Param(r, chi.URLParam(r, "commentId"))
	if err != nil {
		app.badRequestResponse(w, r, err)
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

	// apply incoming fields to the loaded comment
	comment.Content = payLoad.Content
	comment.PostID = payLoad.PostID
	comment.UserID = payLoad.UserID

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

	// validate referenced User exists
	if _, err := app.Store.Users.GetByID(ctx, payLoad.UserID); err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponse(w, r, fmt.Errorf("user_id %d is invalid", payLoad.UserID))
		default:
			app.statusInternalServerError(w, r, err)
		}
		return
	}

	if err := app.Store.Comments.UpdateComment(ctx, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}

	if err := writeJson(w, http.StatusOK, comment); err != nil {
		app.statusInternalServerError(w, r, err)
		return
	}
}
