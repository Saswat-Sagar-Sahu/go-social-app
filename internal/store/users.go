package store

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  string   `json:"-"`
	Activated bool     `json:"activated"`
	CreatedAt string   `json:"created_at"`
	Roles     []string `json:"roles,omitempty"`
}

type UsersStore struct {
	db *sql.DB
}

type UserFilter struct {
	Page     int
	PageSize int
	Name     string
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}
	// assign default role
	if err := s.AssignRole(ctx, user.ID, "user"); err != nil {
		return err
	}
	// populate roles slice
	roles, err := s.GetRoles(ctx, user.ID)
	if err != nil {
		return err
	}
	user.Roles = roles
	return nil
}

func (s *UsersStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	query := `SELECT id, username, email, password, activated, created_at FROM users WHERE id = $1`

	user := &User{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Activated,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password, activated, created_at FROM users WHERE email = $1`
	user := &User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Activated,
		&user.CreatedAt,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UsersStore) List(ctx context.Context, filter UserFilter) ([]*User, int, error) {
	where := ""
	countArgs := make([]any, 0, 1)
	if filter.Name != "" {
		where = " WHERE username ILIKE $1 "
		countArgs = append(countArgs, "%"+filter.Name+"%")
	}

	countQuery := "SELECT COUNT(*) FROM users" + where
	var total int
	if err := s.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	dataQuery := fmt.Sprintf(
		"SELECT id, username, email, activated, created_at FROM users%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		where,
		len(countArgs)+1,
		len(countArgs)+2,
	)

	args := append(countArgs, filter.PageSize, offset)
	rows, err := s.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*User, 0, filter.PageSize)
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Activated,
			&user.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UsersStore) Activate(ctx context.Context, userID int64) error {
	query := `UPDATE users SET activated = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	res, err := s.db.ExecContext(ctx, query, userID)
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

// GetRoles returns role names assigned to a user
func (s *UsersStore) GetRoles(ctx context.Context, userID int64) ([]string, error) {
	query := `SELECT r.name FROM roles r INNER JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = $1`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var roles []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		roles = append(roles, name)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return roles, nil
}

// AssignRole assigns a role to a user by role name. It is idempotent.
func (s *UsersStore) AssignRole(ctx context.Context, userID int64, roleName string) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, (SELECT id FROM roles WHERE name = $2)) ON CONFLICT DO NOTHING`
	_, err := s.db.ExecContext(ctx, query, userID, roleName)
	return err
}

// RemoveRole removes a role assignment from a user.
func (s *UsersStore) RemoveRole(ctx context.Context, userID int64, roleName string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = (SELECT id FROM roles WHERE name = $2)`
	_, err := s.db.ExecContext(ctx, query, userID, roleName)
	return err
}
