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

func (r *Repository) ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error) {
	queryParams, vals := make([]string, len(partNumbers)), make([]interface{}, len(partNumbers))
	for i, v := range partNumbers {
		queryParams[i], vals[i] = fmt.Sprintf("$%d", i+1), v
	}

	query := fmt.Sprintf("select id from products where part_number in (%s)", strings.Join(queryParams, ","))
	rows, err := r.db.QueryContext(ctx, query, vals...)
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

func (r *Repository) WarehousesByProductIDs(ctx context.Context, productIDs []int, lat, long float64) ([]models.WarehouseProduct, error) {
	queryParams, values := make([]string, len(productIDs)), []interface{}{lat, long}
	for i := 0; i < len(productIDs); i++ {
		queryParams[i] = fmt.Sprintf("$%d", i+3)
		values = append(values, productIDs[i])
	}

	query := fmt.Sprintf(
		`SELECT
			wp.product_id,
			wp.warehouse_id,
			ST_Distance(
				ST_Transform(ST_SetSRID(ST_MakePoint($1, $2), 4326), 3857), 
				ST_Transform(ST_SetSRID(ST_MakePoint(w.lat, w.lng), 4326), 3857)
			) AS distance
		FROM warehouse_products wp
		JOIN warehouses w ON wp.warehouse_id = w.id AND available = true
		WHERE wp.product_id IN (%s) AND wp.quantity > 0
		ORDER BY distance`,
		strings.Join(queryParams, ","))

	rows, err := r.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var warehouses []models.WarehouseProduct
	for rows.Next() {
		var wh models.WarehouseProduct
		if err := rows.Scan(
			&wh.ProductID,
			&wh.WarehouseID,
			&wh.Distance,
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

func (r *Repository) SetProductsToReserved(ctx context.Context, reservationID uuid.UUID, warehousesProducts map[int]int) (uuid.UUID, error) {
	var warehouseIDs, productIDs []int
	for pID, wID := range warehousesProducts {
		productIDs = append(productIDs, pID)
		warehouseIDs = append(warehouseIDs, wID)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to begin tx: %w", err)
	}
	defer tx.Rollback()

	warehouseProductIDs, err := r.warehouseProductIDs(ctx, tx, warehouseIDs, productIDs)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to get warehouse product ids: %w", err)
	}

	queryParams, values := make([]string, len(warehouseProductIDs)), make([]interface{}, 0, len(warehouseProductIDs)*3)
	for i, j := 0, 0; i < len(warehouseProductIDs); i, j = i+1, j+3 {
		queryParams[i] = fmt.Sprintf("($%d, $%d, $%d)", j+1, j+2, j+3)
		values = append(values, reservationID, warehouseProductIDs[i], 1)
	}

	query := fmt.Sprintf(
		`INSERT INTO reserved_products (reservation_id, warehouse_product_id, quantity) 
		VALUES %s 
		RETURNING reservation_id`,
		strings.Join(queryParams, ", "))
	row := tx.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return uuid.Nil, fmt.Errorf("error to set products to reserved: %w", row.Err())
	}

	var rID uuid.UUID
	if err = row.Scan(&rID); err != nil {
		return uuid.Nil, fmt.Errorf("scan error: %w", err)
	}

	query = "update warehouse_products set quantity = quantity - 1 where warehouse_id = ANY($1) and product_id = ANY($2)"
	if _, err = tx.ExecContext(ctx, query, pq.Array(warehouseIDs), pq.Array(productIDs)); err != nil {
		return uuid.Nil, fmt.Errorf("error to update quantity in warehouse_products: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return uuid.Nil, fmt.Errorf("commit error: %w", err)
	}
	return rID, nil
}

func (r *Repository) warehouseProductIDs(ctx context.Context, tx *sql.Tx, warehouseIDs, productIDs []int) ([]int, error) {
	query := "select id from warehouse_products where warehouse_id = ANY($1) and product_id = ANY($2)"
	rows, err := tx.QueryContext(ctx, query, pq.Array(warehouseIDs), pq.Array(productIDs))
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

func (r *Repository) Products(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.QueryContext(ctx, "select id, title, part_number, width, height, depth from products")
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

func (r *Repository) AvailabilityProductsByWarehouseID(ctx context.Context, warehouseID int) ([]models.AvailabilityProducts, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			p.id,
			p.part_number,
			p.title,
			wp.quantity,
			w.available
		FROM warehouse_products wp
		JOIN warehouses w ON wp.warehouse_id = w.id
		JOIN products p ON wp.product_id = p.id
		WHERE wp.warehouse_id = $1`,
		warehouseID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var availabilityProducts []models.AvailabilityProducts
	for rows.Next() {
		var availabilityProduct models.AvailabilityProducts
		if err := rows.Scan(
			&availabilityProduct.Product.ID,
			&availabilityProduct.Product.PartNumber,
			&availabilityProduct.Product.Title,
			&availabilityProduct.Quantity,
			&availabilityProduct.WarehouseAvail,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		availabilityProducts = append(availabilityProducts, availabilityProduct)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return availabilityProducts, nil
}

func (r *Repository) SetProductsToConfirmedOrCanceledByProductIDs(ctx context.Context, status int, reservationData models.CancelORConfirmProductsRequest) error {
	productIDs, err := r.ProductsIDsByPartNumbers(ctx, reservationData.PartNumbers)
	if err != nil {
		return fmt.Errorf("error to get products ids: %w", err)
	}

	if len(productIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to begin tx: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(
		ctx,
		`UPDATE reserved_products rp
		SET status = $1
		FROM warehouse_products wp
		JOIN warehouses w ON wp.warehouse_id = w.id
		JOIN products p ON wp.product_id = p.id
		WHERE rp.warehouse_product_id = wp.id AND rp.reservation_id = $2 AND wp.product_id = ANY($3) AND rp.status = 0
		RETURNING rp.warehouse_product_id`,
		status, reservationData.ReservationID, pq.Array(productIDs),
	); 
	if err != nil {
		return fmt.Errorf("error to set products to confirmed: %w", err)
	}

	var warehouseProductIDs []int
	for rows.Next() {
		var warehouseProductID int
		if err := rows.Scan(&warehouseProductID); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		warehouseProductIDs = append(warehouseProductIDs, warehouseProductID)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	if _, err := tx.ExecContext( 
		ctx,
		`UPDATE warehouse_products SET quantity = quantity + 1 WHERE id = ANY($1)`,
		pq.Array(warehouseProductIDs),
	); err != nil {
		return fmt.Errorf("error to update warehouse products: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit tx: %w", err)
	}

	return nil
}

func (r *Repository) SetProductsToConfirmedOrCanceled(ctx context.Context, status int, reservationID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error to begin tx: %w", err)
	}
	defer tx.Rollback()
	
	rows, err := tx.QueryContext(
		ctx,
		`UPDATE reserved_products SET status = $1 WHERE reservation_id = $2 and status = 0 RETURNING warehouse_product_id`,
		status, reservationID)
	if err != nil {
		return fmt.Errorf("error to set products to confirmed: %w", err)
	}

	var warehouseProductIDs []int
	for rows.Next() {
		var warehouseProductID int
		if err := rows.Scan(&warehouseProductID); err != nil {
			return fmt.Errorf("scan error: %w", err)
		}
		warehouseProductIDs = append(warehouseProductIDs, warehouseProductID)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE warehouse_products SET quantity = quantity + 1 WHERE id = ANY($1)`,
		pq.Array(warehouseProductIDs),
	); err != nil {
		return fmt.Errorf("error to set products to confirmed: %w", err)
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error to commit tx: %w", err)
	}

	return nil
}