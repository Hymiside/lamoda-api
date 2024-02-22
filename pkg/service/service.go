package service

import (
	"context"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/Hymiside/lamoda-api/pkg/repository"
)

type store interface {
	ReservationProducts(ctx context.Context, data models.ProductReservationRequest) error
	Warehouses(ctx context.Context) ([]models.Warehouse, error)
	Products(ctx context.Context) ([]models.Product, error)
	WarehouseProducts(ctx context.Context) ([]models.WarehouseProduct, error)
	ReservedProducts(ctx context.Context) ([]models.WarehouseProduct, error)
}

type Service struct {
	S store
}

func NewService(repos *repository.Repository) *Service {
	return &Service{S: newStoreService(repos)}
}