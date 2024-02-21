package handler

import (
	"github.com/Hymiside/lamoda-api/pkg/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	serv *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{serv: serv}
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