package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) ProductsIDsByPartNumbers(ctx context.Context, items []models.ReservationItems) ([]int, error) {
	args, valArgs := make([]string, len(items)), make([]interface{}, len(items))
	for i, v := range items {
		args[i], valArgs[i] = fmt.Sprintf("$%d", i+1), v.PartNumber
	}

	prepareReq := fmt.Sprintf("select id from products where part_number in (%s)", strings.Join(args, ","))
	rows, err := s.db.QueryContext(ctx, prepareReq, valArgs...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return ids, nil
}

func (s *Repository) WarehouseIDsByProductID(ctx context.Context, productIDs []int) ([]int, error) {
	args, valArgs := make([]string, len(productIDs)), make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i], valArgs[i] = fmt.Sprintf("$%d", i+1), v
	}

	prepareReq := fmt.Sprintf("select warehouse_id from warehouse_products where product_id in (%s) and quantity > 0", strings.Join(args, ","))
	rows, err := s.db.QueryContext(ctx, prepareReq, valArgs...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return ids, nil
}

func (s *Repository) WarehousesByIDs(ctx context.Context, warehouseIDs []int) ([]models.Warehouse, error) {
	args, valArgs := make([]string, len(warehouseIDs)), make([]interface{}, len(warehouseIDs))
	for i, v := range warehouseIDs {
		args[i], valArgs[i] = fmt.Sprintf("$%d", i+1), v
	}

	prepareReq := fmt.Sprintf("select id, title, lat, lng from warehouses where id in (%s) and available", strings.Join(args, ","))
	rows, err := s.db.QueryContext(ctx, prepareReq, valArgs...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var warehouse models.Warehouse
		if err := rows.Scan(
			&warehouse.ID,
			&warehouse.Title,
			&warehouse.Latitude,
			&warehouse.Longitude,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		warehouses = append(warehouses, warehouse)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return warehouses, nil
}

func (s *Repository) SetProductsToReserved(ctx context.Context, warehouseID int, productIDs []int) error {
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

func (s *Repository) Products(ctx context.Context) ([]models.Product, error) {
	rows, err := s.db.QueryContext(ctx, "select id, title, part_number, dimensions from products")
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.Title,
			&product.PartNumber,
			&product.Width,
			&product.Height,
			&product.Depth,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return products, nil
}
