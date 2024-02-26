package service

import (
	"context"
	"fmt"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/google/uuid"
)

type repository interface {
	Products(ctx context.Context) ([]models.Product, error)
	ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error)
	AvailabilityProductsByWarehouseID(ctx context.Context, warehouseID int) ([]models.AvailabilityProducts, error)
	WarehousesByProductIDs(ctx context.Context, productIDs []int, lat, long float64) ([]models.WarehouseProduct, error)

	SetProductsToReserved(ctx context.Context, reservationID uuid.UUID, warehousesProducts map[int]int) (uuid.UUID, error)
	SetProductsToConfirmedOrCanceledByProductIDs(ctx context.Context, status int, reservationData models.CancelORConfirmProductsRequest) error
	SetProductsToConfirmedOrCanceled(ctx context.Context, status int, reservationID uuid.UUID) error
}

type Service struct {
	repos repository
}

func NewService(repos repository) *Service {
	return &Service{repos: repos}
}

func (s *Service) ReservationProducts(ctx context.Context, req models.ReservationProductsRequest) (uuid.UUID, error) {
	productIDs, err := s.repos.ProductsIDsByPartNumbers(ctx, req.PartNumbers)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to get products: %v", err)
	}

	if len(productIDs) == 0 {
		return uuid.Nil, fmt.Errorf("products not found")
	}

	warehousesProducts, err := s.repos.WarehousesByProductIDs(ctx, productIDs, req.Latitude, req.Longitude)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to get warehouses: %v", err)
	}

	if len(warehousesProducts) == 0 {
		return uuid.Nil, fmt.Errorf("warehousesProducts not found")
	}

	reservation := make(map[int]int)
	for _, v := range warehousesProducts {
		if _, ok := reservation[v.ProductID]; !ok {
			reservation[v.ProductID] = v.WarehouseID
		}
	}

	reservationID, err := s.repos.SetProductsToReserved(ctx, uuid.New(), reservation)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error to set products to reserved: %v", err)
	}
	return reservationID, nil
}

func (s *Service) Products(ctx context.Context) ([]models.Product, error) {
	products, err := s.repos.Products(ctx)
	if err != nil {
		return nil, fmt.Errorf("error to get products: %v", err)
	}
	return products, nil
}

func (s *Service) AvailabilityProductsByWarehouseID(ctx context.Context, warehouseID int) ([]models.AvailabilityProducts, error) {
	reservedProducts, err := s.repos.AvailabilityProductsByWarehouseID(ctx, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("error to get reserved products: %v", err)
	}
	return reservedProducts, nil
}

func (s *Service) ConfirmOrCancelReservedProducts(ctx context.Context, status int, req models.CancelORConfirmProductsRequest) error {

	if req.PartNumbers == nil {
		if err := s.repos.SetProductsToConfirmedOrCanceled(ctx, status, req.ReservationID); err != nil {
			return fmt.Errorf("error to set products to confirmed: %v", err)
		}
		return nil
	}

	if err := s.repos.SetProductsToConfirmedOrCanceledByProductIDs(ctx, status, req); err != nil {
		return fmt.Errorf("error to set products to confirmed: %v", err)
	}
	return nil
}
