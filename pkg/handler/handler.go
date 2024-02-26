package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

//go:generate mockery --name=service --output=../../mock/service --outpkg=service_mock --filename=service_mock.go
type service interface {
	Products(ctx context.Context) ([]models.Product, error)
	AvailabilityProductsByWarehouseID(ctx context.Context, warehouseID int) ([]models.AvailabilityProducts, error)

	ReservationProducts(ctx context.Context, data models.ReservationProductsRequest) (uuid.UUID, error)
	ConfirmOrCancelReservedProducts(ctx context.Context, status int, req models.CancelORConfirmProductsRequest) error
}

type Handler struct {
	services service
	validate *validator.Validate
}

func NewHandler(service service) *Handler {
	return &Handler{
		services: service,
		validate: validator.New(),
	}
}

func (h *Handler) NewRoutes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/products", h.products)
	mux.Get("/products/availability", h.availabilityProduct)
	mux.Post("/reservation-products", h.reservationProducts)
	mux.Delete("/reservation-products", h.cancelReservationProducts)
	mux.Post("/confirm-reservation", h.confirmReservationProducts)

	return mux
}

func (h *Handler) products(w http.ResponseWriter, r *http.Request) {
	products, err := h.services.Products(r.Context())
	if err != nil {
		log.Errorf("error to get products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Errorf("error to encode products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) availabilityProduct(w http.ResponseWriter, r *http.Request) {
	queryVal := r.URL.Query().Get("warehouse_id")
	if queryVal == "" {
		http.Error(w, "warehouse_id is required", http.StatusBadRequest)
		return
	}

	warehouseID, err := strconv.Atoi(queryVal)
	if err != nil {
		log.Errorf("error to convert warehouse_id: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reservedProducts, err := h.services.AvailabilityProductsByWarehouseID(r.Context(), warehouseID)
	if err != nil {
		log.Errorf("error to get reserved products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(reservedProducts); err != nil {
		log.Errorf("error to encode products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) reservationProducts(w http.ResponseWriter, r *http.Request) {
	var req models.ReservationProductsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("error to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.validate.Struct(req); err != nil {
		log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	reservationID, err := h.services.ReservationProducts(r.Context(), req)
	if err != nil {
		log.Errorf("error to reservation products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(
		map[string]uuid.UUID{"reservation_id": reservationID},
	); err != nil {
		log.Errorf("error to encode products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) cancelReservationProducts(w http.ResponseWriter, r *http.Request) {
	var req models.CancelORConfirmProductsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("error to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.validate.Struct(req); err != nil {
		log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.services.ConfirmOrCancelReservedProducts(r.Context(), 2, req); err != nil {
		log.Errorf("error to cancel reserved products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) confirmReservationProducts(w http.ResponseWriter, r *http.Request) {
	var req models.CancelORConfirmProductsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("error to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = h.validate.Struct(req); err != nil {
		log.Errorf("validation error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.services.ConfirmOrCancelReservedProducts(r.Context(), 2, req); err != nil {
		log.Errorf("error to confirm reserved products: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
