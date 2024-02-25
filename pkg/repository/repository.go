package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (s *Repository) ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error) {
	queryParams, vals := make([]string, len(partNumbers)), make([]interface{}, len(partNumbers))
	for i, v := range partNumbers {
		queryParams[i], vals[i] = fmt.Sprintf("$%d", i+1), v
	}

	query := fmt.Sprintf("select id from products where part_number in (%s)", strings.Join(queryParams, ","))
	rows, err := s.db.QueryContext(ctx, query, vals...)
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

func (s *Repository) WarehousesByProductIDs(ctx context.Context, productIDs []int) ([]models.WarehouseProductID, error) {
	queryParams, values := make([]string, len(productIDs)), make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		queryParams[i], values[i] = fmt.Sprintf("$%d", i+1), v
	}

	query := fmt.Sprintf(
		`select
			wp.product_id,
			wp.quantity,
			w.id,
			w.lat,
			w.lng
		from warehouse_products wp 
		join warehouses w on 
			wp.warehouse_id = w.id and w.available = true
		where wp.product_id in (%s) and wp.quantity > 0`,
		strings.Join(queryParams, ","))

	rows, err := s.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var warehouses []models.WarehouseProductID
	for rows.Next() {
		var wh models.WarehouseProductID
		if err := rows.Scan(
			&wh.ProductID,
			&wh.Quantity,
			&wh.Warehouse.ID,
			&wh.Warehouse.Latitude,
			&wh.Warehouse.Longitude,
		); err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		warehouses = append(warehouses, wh)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}
	return warehouses, nil
}

func (s *Repository) SetProductsToReserved(ctx context.Context, reservationID uuid.UUID, warehousesProducts []models.ReservationProducts) error {
	var warehouseIDs, productIDs []int
	for _, v := range warehousesProducts {
		warehouseIDs = append(warehouseIDs, v.WarehouseID)
		productIDs = append(productIDs, v.ProductID)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	warehouseProductIDs, err := s.warehouseProductIDs(ctx, tx, warehouseIDs, productIDs)
	if err != nil {
		return fmt.Errorf("error to get warehouse product ids: %w", err)
	}

	queryParams, values := make([]string, len(warehouseProductIDs)), make([]interface{}, 0, len(warehouseProductIDs)*3)
	for i, j := 0, 0; i < len(warehouseProductIDs); i, j = i+1, j+3 {
		queryParams[i] = fmt.Sprintf("($%d, $%d, $%d)", j+1, j+2, j+3)
		values = append(values, reservationID, warehouseProductIDs[i], 1)
	}

	query := fmt.Sprintf("insert into reserved_products (reservation_id, warehouse_product_id, quantity) values %s", strings.Join(queryParams, ", "))
	if _, err = tx.ExecContext(ctx, query, values...); err != nil {
		return fmt.Errorf("error to set products to reserved: %w", err)
	}

	query = "update warehouse_products set quantity = quantity - 1 where warehouse_id = ANY($1) and product_id = ANY($2)"
	if _, err = tx.ExecContext(ctx, query, pq.Array(warehouseIDs), pq.Array(productIDs)); err != nil {
		return fmt.Errorf("error to update quantity in warehouse_products: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}

func (s *Repository) warehouseProductIDs(ctx context.Context, tx *sql.Tx, warehouseIDs, productIDs []int) ([]int, error) {
	query := "select id from warehouse_products where warehouse_id = ANY($1) and product_id = ANY($2)"
	rows, err := tx.QueryContext(ctx, query,pq.Array(warehouseIDs), pq.Array(productIDs))
	if err != nil {
		return nil, fmt.Errorf("error to get warehousesids: %w", err)
	}
	defer rows.Close()

	var warehouseProductIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		warehouseProductIDs = append(warehouseProductIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return warehouseProductIDs, nil
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
