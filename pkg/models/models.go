package models

import "github.com/google/uuid"

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

type ReservationProductsRequest struct {
	PartNumbers []string `json:"part_numbers" validate:"required,min=1"`
	Latitude    float64  `json:"latitude" validate:"required"`
	Longitude   float64  `json:"longitude" validate:"required"`
}

type CancelORConfirmProductsRequest struct {
	ReservationID uuid.UUID `json:"reservation_id" validate:"required"`
	PartNumbers   []string  `json:"part_numbers"`
}

type Warehouse struct {
	ID        int     `json:"id"`
	Title     string  `json:"title"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"long"`
}

type Product struct {
	ID         int    `json:"id"`
	PartNumber string `json:"part_number"`
	Title      string `json:"title"`
	Width      int    `json:"width,omitempty"`
	Height     int    `json:"height,omitempty"`
	Depth      int    `json:"depth,omitempty"`
}

type WarehouseProduct struct {
	ProductID   int
	WarehouseID int
	Distance    float64
}

type ReservationProducts struct {
	ProductID   int
	WarehouseID int
}

type AvailabilityProducts struct {
	Product        Product
	WarehouseAvail bool
	Quantity       int
}
