package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hymiside/lamoda-api/pkg/models"
	_ "github.com/lib/pq"
)

type store interface{
	ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error)
	WarehouseIDsByProductID(ctx context.Context, productIDs []int) ([]int, error)
	WarehousesByIDs(ctx context.Context, warehouseIDs []int) ([]models.Warehouse, error)
	SetProductsToReserved(ctx context.Context, warehouseID int, productIDs []int) error
}

type Repository struct {
	S store
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{S: newStoreRepository(db)}
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
