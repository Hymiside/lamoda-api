package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/lib/pq"
)

type storeRepository struct {
	db *sql.DB
}

func newStoreRepository(db *sql.DB) *storeRepository {
	return &storeRepository{db: db}
}

func (s *storeRepository) ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error) {
	rows, err := s.db.QueryContext(ctx, "select id from products where part_number = ANY($1)", pq.Array(partNumbers))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *storeRepository) WarehouseIDsByProductID(ctx context.Context, productIDs []int) ([]int, error) {
	rows, err := s.db.QueryContext(ctx, "select warehouse_id from warehouse_products where product_id = ANY($1) and quantity > 0", pq.Array(productIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *storeRepository) WarehousesByIDs(ctx context.Context, warehouseIDs []int) ([]models.Warehouse, error) {
	rows, err := s.db.QueryContext(ctx, "select id, title, lat, lng from warehouses where id = ANY($1) and available = true", pq.Array(warehouseIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var warehouse models.Warehouse
		if err := rows.Scan(&warehouse.ID, &warehouse.Title, &warehouse.Lat, &warehouse.Long); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, warehouse)
	}
	return warehouses, nil
}

func (s *storeRepository) SetProductsToReserved(ctx context.Context, warehouseID int, productIDs []int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	rows, err := tx.QueryContext(ctx, "select id from warehouse_products where warehouse_id = $1 and product_id = ANY($2)", warehouseID, pq.Array(productIDs))
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	var warehouseProductIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return err
		}
		warehouseProductIDs = append(warehouseProductIDs, id)
	}

	vals := []interface{}{}
	for _, id := range warehouseProductIDs {
		vals = append(vals, id, 1)
	}
	
	_, err = tx.ExecContext(ctx, "insert into reserved_products (warehouse_product_id, quantity) values ($1, $2)", vals...)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error to set products to reserved: %w", err)
	}

	_, err = tx.ExecContext(
		ctx, 
		"update warehouse_products set quantity = quantity - 1 where warehouse_id = $1 and product_id = ANY($2)", 
		warehouseID, pq.Array(productIDs),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}