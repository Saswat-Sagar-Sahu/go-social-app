package store

import (
	"context"
	"database/sql"
	"errors"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAt  string `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

var ErrResourceExists = errors.New("resource already exists")

func (s *FollowersStore) AddFollower(ctx context.Context, userID, followerID int64) error {

	checkQuery := `SELECT 1 FROM followers WHERE user_id = $1 AND follower_id = $2`
	var exists int
	err := s.db.QueryRowContext(ctx, checkQuery, userID, followerID).Scan(&exists)
	if err == nil {
		return ErrResourceExists
	}
	if err != sql.ErrNoRows {
		return err
	}

	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	_, err = s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}
	return nil
}

func (s *FollowersStore) RemoveFollower(ctx context.Context, userID, followerID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}
	return nil
}
