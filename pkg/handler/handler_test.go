package handler_test

import (
    "bytes"
	"encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/Hymiside/lamoda-api/pkg/handler"
    "github.com/Hymiside/lamoda-api/pkg/models"
    mockservice "github.com/Hymiside/lamoda-api/mock/service"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestHandler_products(t *testing.T) {
    svc := new(mockservice.ServiceMock)
    h := handler.NewHandler(svc)

    svc.On("Products", mock.Anything).Return([]models.Product{}, nil)

    req, err := http.NewRequest("GET", "/products", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    h.NewRoutes().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    svc.AssertExpectations(t)
}

func TestHandler_availabilityProduct(t *testing.T) {
    svc := new(mockservice.ServiceMock)
    h := handler.NewHandler(svc)

    warehouseID := 1
    svc.On("AvailabilityProductsByWarehouseID", mock.Anything, warehouseID).Return([]models.AvailabilityProducts{}, nil)

    req, err := http.NewRequest("GET", "/products/availability?warehouse_id=1", nil)
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    h.NewRoutes().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    svc.AssertExpectations(t)
}

func TestHandler_reservationProducts(t *testing.T) {
    svc := new(mockservice.ServiceMock)
    h := handler.NewHandler(svc)

    reservationRequest := models.ReservationProductsRequest{
        PartNumbers: []string{"P13579", "P97431", "P13279"},
        Latitude:    21.213,
        Longitude:   32.23,
    }
    reservationID := uuid.New()
    svc.On("ReservationProducts", mock.Anything, reservationRequest).Return(reservationID, nil)

    requestBody, _ := json.Marshal(reservationRequest)
    req, err := http.NewRequest("POST", "/reservation-products", bytes.NewBuffer(requestBody))
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    h.NewRoutes().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    svc.AssertExpectations(t)
}

func TestHandler_cancelReservationProducts(t *testing.T) {
    svc := new(mockservice.ServiceMock)
    h := handler.NewHandler(svc)

    cancelRequest := models.CancelORConfirmProductsRequest{
        ReservationID: uuid.New(),
        PartNumbers:   []string{"P13579"},
    }
    svc.On("ConfirmOrCancelReservedProducts", mock.Anything, 2, cancelRequest).Return(nil)

    requestBody, _ := json.Marshal(cancelRequest)
    req, err := http.NewRequest("DELETE", "/reservation-products", bytes.NewBuffer(requestBody))
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    h.NewRoutes().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    svc.AssertExpectations(t)
}

func TestHandler_confirmReservationProducts(t *testing.T) {
    svc := new(mockservice.ServiceMock)
    h := handler.NewHandler(svc)

    confirmRequest := models.CancelORConfirmProductsRequest{
        ReservationID: uuid.New(),
        PartNumbers:   []string{"P13579"},
    }
    svc.On("ConfirmOrCancelReservedProducts", mock.Anything, 2, confirmRequest).Return(nil)

    requestBody, _ := json.Marshal(confirmRequest)
    req, err := http.NewRequest("POST", "/confirm-reservation", bytes.NewBuffer(requestBody))
    assert.NoError(t, err)

    rr := httptest.NewRecorder()
    h.NewRoutes().ServeHTTP(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
    svc.AssertExpectations(t)
}