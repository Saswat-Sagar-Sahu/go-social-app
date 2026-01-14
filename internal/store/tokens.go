package store

import (
	"context"
	"database/sql"
	"time"
)

type Token struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt string    `json:"created_at"`
}

type TokenStore struct {
	db *sql.DB
}

func (s *TokenStore) Create(ctx context.Context, t *Token) error {
	query := `INSERT INTO user_tokens (user_id, token, expires_at) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query, t.UserID, t.Token, t.ExpiresAt).Scan(&t.ID, &t.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *TokenStore) GetByToken(ctx context.Context, token string) (*Token, error) {
	query := `SELECT id, user_id, token, expires_at, created_at FROM user_tokens WHERE token = $1`
	t := &Token{}
	err := s.db.QueryRowContext(ctx, query, token).Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.CreatedAt)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return t, nil
}

func (s *TokenStore) DeleteByToken(ctx context.Context, token string) error {
	query := `DELETE FROM user_tokens WHERE token = $1`
	_, err := s.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	return nil
}
