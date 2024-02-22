package models

import "time"

type ConfigPostgres struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type ConfigServer struct {
	Host string
	Port string
}

type ProductReservationRequest struct {
	PartNumbers []string `json:"part_numbers" validate:"required,min=1"`
	Quantity    int      `json:"quantity" validate:"required,min=1"`
	Lat         float64  `json:"lat" validate:"required"`
	Long        float64  `json:"long" validate:"required"`
}

type Warehouse struct {
	ID        int `db:"id" json:"id"`
	Title     string `db:"title" json:"title"`
	Lat       float64 `db:"lat" json:"lat"`
	Long      float64 `db:"lng" json:"long"`
}

type Product struct {
	ID          int `db:"id" json:"id"` 
	PartNumber  string `db:"part_number" json:"part_number"`
	Title       string `db:"title" json:"title"`
	Dimensions  *map[string]int `db:"dimensions" json:"dimensions"`
}

type WarehouseProduct struct {
	Product Product `json:"product"`
	Warehouse Warehouse `json:"warehouse"`
	Quantity int `json:"quantity"`
}

type ReservedProducts struct {
	Product Product `json:"product"`
	Warehouse Warehouse `json:"warehouse"`
	Quantity int `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}