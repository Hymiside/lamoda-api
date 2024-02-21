package repository

import "database/sql"

type storeRepository struct {
	db *sql.DB
}

func newStoreRepository(db *sql.DB) *storeRepository {
	return &storeRepository{db: db}
}