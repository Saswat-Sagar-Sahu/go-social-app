package auth

import "context"

type ctxKey string

const ContextUserKey ctxKey = "userID"

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(ContextUserKey)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

type Authenticator interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, password string) error
	GenerateToken(ctx context.Context, userID int64) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}
