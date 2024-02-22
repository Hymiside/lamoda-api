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

// products, warehoses, reservation-products, reserved-products, cancel-reservation-products, buy-reserved-products

func (h *Handler) NewRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/products", h.products)
	mux.Get("/warehoses", h.warehoses)
	mux.Post("/reservation-products", h.reservationProducts)
	mux.Post("/reserved-products", h.reservedProducts)
	mux.Post("/cancel-reservation-products", h.cancelReservationProducts)
	mux.Post("/buy-reserved-products", h.buyReservedProducts)

	return mux
}
