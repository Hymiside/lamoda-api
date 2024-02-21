package repository

import (
	"context"
	"fmt"

	"database/sql"

	"github.com/Hymiside/lamoda-api/pkg/models"
)

type store interface{}

type Repository struct {
	s store
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{s: newStoreRepository(db)}
}

func NewPostgresDB(ctx context.Context, c models.ConfigPostgres) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.Name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("error to connection postgres: %v", err)
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		db.Close()
	}(ctx)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("connection test error: %w", err)
	}

	return db, nil
}
