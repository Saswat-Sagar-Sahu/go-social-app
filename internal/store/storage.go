package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		GetUserFeed(context.Context, int64) ([]*Post, error)
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostId(context.Context, int64) (error, *Comment)
		GetByUserId(context.Context, int64) (error, *Comment)
		GetByCommentId(context.Context, int64) (error, *Comment)
		DeleteByID(context.Context, int64) error
		UpdateComment(context.Context, *Comment) error
	}
	Followers interface {
		AddFollower(context.Context, int64, int64) error
		RemoveFollower(context.Context, int64, int64) error
	}
	Tokens interface {
		Create(context.Context, *Token) error
		GetByToken(context.Context, string) (*Token, error)
		DeleteByToken(context.Context, string) error
	}
	Invitations interface {
		Create(context.Context, *Invitation) error
		GetByToken(context.Context, string) (*Invitation, error)
		MarkAccepted(context.Context, string, int64) error
	}
}

func (s Storage) Post() {
	panic("unimplemented")
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:       &PostStore{db: db},
		Users:       &UsersStore{db: db},
		Comments:    &CommentStore{db: db},
		Followers:   &FollowersStore{db: db},
		Tokens:      &TokenStore{db: db},
		Invitations: &InvitationsStore{db: db},
	}
}
