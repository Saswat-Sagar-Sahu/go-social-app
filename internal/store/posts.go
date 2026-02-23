package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64    `json:"id"`
	Content   string   `json:"content"`
	Title     string   `json:"title"`
	UserID    *int64   `json:"user_id,omitempty"`
	ImageURL  *string  `json:"image_url,omitempty"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// ErrNotFound is returned when a requested entity is not found in the store.
var ErrNotFound = errors.New("No record found")

type PostStore struct {
	db *sql.DB
}

func (s *PostStore) Create(ctx context.Context, post *Post) error {
	if post.UserID == nil || *post.UserID == 0 {
		return errors.New("post user_id is required")
	}
	userID := *post.UserID

	// cast $3 to bigint to ensure the DB receives an integer type
	query := `INSERT INTO posts (content, title, user_id, image_url, tags) VALUES ($1, $2, $3::bigint, $4, $5) RETURNING id, created_at, updated_at`
	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		userID,
		post.ImageURL,
		pq.Array(post.Tags),
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `SELECT id, content, title, user_id, image_url, tags, created_at, updated_at FROM posts WHERE id = $1`

	post := &Post{}
	var userID sql.NullInt64
	var imageURL sql.NullString
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&userID,
		&imageURL,
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
	if userID.Valid {
		v := userID.Int64
		post.UserID = &v
	} else {
		post.UserID = nil
	}
	if imageURL.Valid {
		v := imageURL.String
		post.ImageURL = &v
	} else {
		post.ImageURL = nil
	}

	return post, nil
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]*Post, error) {
	// query := `SELECT p.id, p.content, p.title, p.user_id, p.tags, p.created_at, p.updated_at
	// FROM posts p
	// JOIN followers f ON p.user_id = f.user_id
	// WHERE f.follower_id = $1
	// ORDER BY p.created_at DESC`
	query := `
    SELECT
        p.id, p.content, p.title, p.user_id, p.image_url, p.tags, p.created_at, p.updated_at,
        COALESCE(MAX(c.created_at), p.created_at) AS last_activity
    FROM posts p
    LEFT JOIN followers f ON p.user_id = f.user_id AND f.follower_id = $1
    LEFT JOIN comments c ON c.post_id = p.id
        AND c.user_id IN (SELECT user_id FROM followers WHERE follower_id = $1)
    WHERE f.follower_id = $1 OR c.user_id IS NOT NULL
    GROUP BY p.id, p.content, p.title, p.user_id, p.image_url, p.tags, p.created_at, p.updated_at
    ORDER BY last_activity DESC
    `

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		post := &Post{}
		var lastActivity sql.NullString
		var userID sql.NullInt64
		var imageURL sql.NullString
		err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.Title,
			&userID,
			&imageURL,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
			&lastActivity,
		)
		if err != nil {
			return nil, err
		}
		if userID.Valid {
			v := userID.Int64
			post.UserID = &v
		} else {
			post.UserID = nil
		}
		if imageURL.Valid {
			v := imageURL.String
			post.ImageURL = &v
		} else {
			post.ImageURL = nil
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostStore) GetByUserID(ctx context.Context, userID int64) ([]*Post, error) {
	query := `SELECT id, content, title, user_id, image_url, tags, created_at, updated_at
	          FROM posts
	          WHERE user_id = $1
	          ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]*Post, 0)
	for rows.Next() {
		post := &Post{}
		var ownerID sql.NullInt64
		var imageURL sql.NullString
		if err := rows.Scan(
			&post.ID,
			&post.Content,
			&post.Title,
			&ownerID,
			&imageURL,
			pq.Array(&post.Tags),
			&post.CreatedAt,
			&post.UpdatedAt,
		); err != nil {
			return nil, err
		}
		if ownerID.Valid {
			v := ownerID.Int64
			post.UserID = &v
		}
		if imageURL.Valid {
			v := imageURL.String
			post.ImageURL = &v
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts
	          SET title = $1, content = $2, image_url = $3, tags = $4, updated_at = CURRENT_TIMESTAMP
	          WHERE id = $5
	          RETURNING updated_at`

	if err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ImageURL, pq.Array(post.Tags), post.ID).Scan(&post.UpdatedAt); err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

func (s *PostStore) DeleteByID(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = $1`

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
