package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewPostgres(dsn string) (*sqlx.DB, error) {
	database, err := sqlx.ConnectContext(
		context.Background(),
		"pgx",
		dsn,
	)
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(20)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(30 * time.Minute)

	return database, nil
}
