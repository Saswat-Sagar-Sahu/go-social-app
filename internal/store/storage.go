package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, int64) (error, *Comment)
		GetByUserId(context.Context, int64) (error, *Comment)
		GetByCommentId(context.Context, int64) (error, *Comment)
		DeleteByID(context.Context, int64) error
		UpdateComment(context.Context, *Comment) error
	}
}

func (s Storage) Post() {
	panic("unimplemented")
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UsersStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
