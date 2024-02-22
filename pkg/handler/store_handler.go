package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Hymiside/lamoda-api/pkg/models"
)

func (h *Handler) products(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) warehoses(w http.ResponseWriter, r *http.Request) {}

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
	sendJSONResponse(w, http.StatusOK, nil)
}

func (h *Handler) reservedProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) cancelReservationProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) buyReservedProducts(w http.ResponseWriter, r *http.Request) {}