package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// норм юзать интерфейсы, убрать wg и заюзать структуру, переделать БД, убрать таймаут контекст, убрать ненужные функции, переделать возвраты ошибок
// переделать запросы к БД, написать тесты, проверка rows на ошибку

type service interface {
	ReservationProducts(ctx context.Context, data models.ReservationProductsRequest) error
	Products(ctx context.Context) ([]models.Product, error)
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
	mux.Post("/reservation-products", h.reservationProducts)
	mux.Delete("/reservation-products", h.cancelReservationProducts)
	mux.Post("/confirm-reservation", h.confirmReservationProducts)

	return mux
}

func (h *Handler) products(w http.ResponseWriter, r *http.Request) {
	products, err := h.services.Products(r.Context())
	if err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string][]models.Product{"message": products})
}

func (h *Handler) reservationProducts(w http.ResponseWriter, r *http.Request) {
	var req models.ReservationProductsRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorf("error to decode request: %v", err)
		if err = sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError); err != nil {
			log.Errorf("error to send response: %v", err)
		}
		return
	}

	if err = h.validate.Struct(req); err != nil {
		log.Errorf("validation error: %v", err)
		if err = sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError); err != nil {
			log.Errorf("error to send response: %v", err)
		}
		return
	}

	if err = h.services.ReservationProducts(r.Context(), req); err != nil {
		log.Errorf("error to reservation products: %v", err)
		if err = sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError); err != nil {
			log.Errorf("error to send response: %v", err)
		}
		return
	}
	if err = sendJSONResponse(w, http.StatusOK, map[string]string{"message": "success"}); err != nil {
		log.Errorf("error to send response: %v", err)
	}
}

func (h *Handler) cancelReservationProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) confirmReservationProducts(w http.ResponseWriter, r *http.Request) {}
