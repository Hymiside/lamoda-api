package handler

import (
	"github.com/Hymiside/lamoda-api/pkg/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	services *service.Service
	validate *validator.Validate
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{
		services: serv,
		validate: validator.New(),
	}
}

func (h *Handler) NewRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/products", h.products)
	mux.Get("/warehouses", h.warehouses)
	mux.Get("/reserved-products", h.reservedProducts)
	mux.Get("/warehouse-products", h.warehouseProducts)
	mux.Post("/reservation-products", h.reservationProducts)
	mux.Post("/cancel-reservation-products", h.cancelReservationProducts)
	mux.Post("/buy-reserved-products", h.buyReservedProducts)

	return mux
}
