package handler

import "net/http"

func (h *Handler) products(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) warehoses(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) reservationProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) reservedProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) cancelReservationProducts(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) buyReservedProducts(w http.ResponseWriter, r *http.Request) {}