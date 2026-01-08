package store

import (
	"context"
	cryptorand "crypto/rand"
	"database/sql"
	"errors"
	"log"
	"math/big"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64         `json:"id"`
	Content   string        `json:"content"`
	Title     string        `json:"title"`
	UserID    sql.NullInt64 `json:"user_id"`
	Tags      []string      `json:"tags"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}

// ErrNotFound is returned when a requested entity is not found in the store.
var ErrNotFound = errors.New("No record found")

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	var userID int64
	if post.UserID.Valid && post.UserID.Int64 != 0 {
		userID = post.UserID.Int64
	} else {
		// pick a random existing user id between 1 and 5
		userID = randomInt(1, 5)
		post.UserID = sql.NullInt64{Int64: userID, Valid: true}
	}

	log.Printf("User Id : %d", userID)

	// cast $3 to bigint to ensure the DB receives an integer type
	query := `INSERT INTO posts (content, title, user_id, tags) VALUES ($1, $2, $3::bigint, $4) RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		userID,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func randomInt(min, max int) int64 {
	if min >= max {
		return int64(min)
	}
	span := int64(max - min + 1)
	nBig, err := cryptorand.Int(cryptorand.Reader, big.NewInt(span))
	if err != nil {
		return int64(min)
	}
	return int64(min) + nBig.Int64()
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, content, title, user_id, tags, created_at, updated_at FROM posts WHERE id = $1`

	post := &Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}
