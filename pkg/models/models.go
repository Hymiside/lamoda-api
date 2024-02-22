package models

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
	ID        int
	Title     string
	Lat       float64
	Long      float64
}