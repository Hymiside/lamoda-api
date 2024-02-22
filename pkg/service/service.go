package service

import (
	"context"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/Hymiside/lamoda-api/pkg/repository"
)

type store interface {
	ReservationProducts(_ context.Context, data models.ProductReservationRequest) error
	Warehouses(_ context.Context) ([]models.Warehouse, error)
	Products(_ context.Context) ([]models.Product, error)
	WarehouseProducts(_ context.Context) ([]models.WarehouseProduct, error)
	ReservedProducts(_ context.Context) ([]models.WarehouseProduct, error)
}

type Service struct {
	S store
}

func NewService(repos *repository.Repository) *Service {
	return &Service{S: newStoreService(repos)}
}