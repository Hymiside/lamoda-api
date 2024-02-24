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

type ReservationItems struct {
	PartNumber string `json:"part_number" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required,min=1"`
}

type ReservationProductsRequest struct {
	Items []ReservationItems `json:"items" validate:"required"`
	Latitude      float64 `json:"latitude" validate:"required"`
	Longitude     float64 `json:"longitude" validate:"required"`
}

type Warehouse struct {
	ID    int     `db:"id" json:"id"`
	Title string  `db:"title" json:"title"`
	Latitude   float64 `db:"lat" json:"lat"`
	Longitude  float64 `db:"lng" json:"long"`
}

type Product struct {
	ID         int    `db:"id" json:"id"`
	PartNumber string `db:"part_number" json:"part_number"`
	Title      string `db:"title" json:"title"`
	Width      int    `db:"width" json:"width"`
	Height     int    `db:"height" json:"height"`
	Depth      int    `db:"depth" json:"depth"`
}

type WarehouseProduct struct {
	ProductID   int   `json:"product_id"`
	Warehouse Warehouse `json:"warehouse"`
	Quantity  int       `json:"quantity"`
}

type ReservedProducts struct {
	ResevationID int    `json:"reservation_id"`
	Product   Product   `json:"product"`
	Warehouse Warehouse `json:"warehouse"`
	CreatedAt time.Time `json:"created_at"`
}
