package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConns, MaxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db , err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(duration)
	db.SetMaxIdleConns(MaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}