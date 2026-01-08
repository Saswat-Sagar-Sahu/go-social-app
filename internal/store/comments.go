package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64 `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content).
		Scan(
			&comment.ID,
			&comment.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) GetByPostId(ctx context.Context, postID int64) (error, *Comment) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = $1 order by created_at desc`

	comment := &Comment{}
	err := s.db.QueryRowContext(ctx, query, postID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound, nil
		default:
			return err, nil
		}
	}
	return nil, comment
}

func (s *CommentStore) GetByUserId(ctx context.Context, userID int64) (error, *Comment) {
	query := `SELECT id, post_id, user_id, content, created_at FROM comments WHERE user_id = $1 order by created_at desc`

	comment := &Comment{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserID,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound, nil
		default:
			return err, nil
		}
	}
	return nil, comment
}