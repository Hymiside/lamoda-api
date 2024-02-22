package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Hymiside/lamoda-api/pkg/models"
)

func (h *Handler) products(w http.ResponseWriter, r *http.Request) {
	products, err := h.services.S.Products(r.Context())
	if err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string][]models.Product{"message": products})
}

func (h *Handler) warehouses(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.services.S.Warehouses(r.Context())
	if err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string][]models.Warehouse{"message": warehouses})
}

func (h *Handler) warehouseProducts(w http.ResponseWriter, r *http.Request) {
	warehouseProducts, err := h.services.S.WarehouseProducts(r.Context())
	if err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string][]models.WarehouseProduct{"message": warehouseProducts})
}

func (h *Handler) reservationProducts(w http.ResponseWriter, r *http.Request) {
	var req models.ProductReservationRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.validate.Struct(req); err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.services.S.ReservationProducts(r.Context(), req); err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string]string{"message": "success"})
}

func (h *Handler) reservedProducts(w http.ResponseWriter, r *http.Request) {
	warehouseProducts, err := h.services.S.ReservedProducts(r.Context())
	if err != nil {
		sendJSONErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, http.StatusOK, map[string][]models.WarehouseProduct{"message": warehouseProducts})
}

func (h *Handler) cancelReservationProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) buyReservedProducts(w http.ResponseWriter, r *http.Request) {}