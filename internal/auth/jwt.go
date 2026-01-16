package auth

import (
    "context"
    "errors"
    "time"

    "github.com/Saswat-Sagar-Sahu/Social/internal/store"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type JWTAuth struct {
    secret        []byte
    expiryMinutes int
    store         store.Storage
}

func NewJWTAuth(secret string, expiryMinutes int, s store.Storage) *JWTAuth {
    return &JWTAuth{secret: []byte(secret), expiryMinutes: expiryMinutes, store: s}
}

func (j *JWTAuth) HashPassword(password string) (string, error) {
    h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(h), nil
}

func (j *JWTAuth) ComparePassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func (j *JWTAuth) GenerateToken(ctx context.Context, userID int64) (string, error) {
    now := time.Now()
    exp := now.Add(time.Duration(j.expiryMinutes) * time.Minute)

    claims := jwt.MapClaims{
        "sub": userID,
        "exp": exp.Unix(),
        "iat": now.Unix(),
    }
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := tok.SignedString(j.secret)
    if err != nil {
        return "", err
    }

    // persist token
    t := &store.Token{UserID: userID, Token: signed, ExpiresAt: exp}
    if err := j.store.Tokens.Create(ctx, t); err != nil {
        return "", err
    }

    return signed, nil
}

func (j *JWTAuth) ParseToken(ctx context.Context, tokenStr string) (int64, error) {
    if tokenStr == "" {
        return 0, errors.New("empty token")
    }
    tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return j.secret, nil
    })
    if err != nil {
        return 0, err
    }
    if !tok.Valid {
        return 0, errors.New("invalid token")
    }
    claims, ok := tok.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid claims")
    }
    sub, ok := claims["sub"].(float64)
    if !ok {
        return 0, errors.New("invalid subject claim")
    }

    // check persisted token exists and not expired
    t, err := j.store.Tokens.GetByToken(ctx, tokenStr)
    if err != nil {
        return 0, err
    }
    if time.Now().After(t.ExpiresAt) {
        return 0, errors.New("token expired")
    }

    return int64(sub), nil
}
