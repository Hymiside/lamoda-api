package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/google/uuid"
)

type repository interface {
	Products(ctx context.Context) ([]models.Product, error)
	ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error)
	WarehousesByProductIDs(ctx context.Context, productIDs []int, lat, long float64) ([]models.WarehouseProductID, error)
	SetProductsToReserved(ctx context.Context, reservationID uuid.UUID, warehousesProducts map[int]int) error
}

type Service struct {
	repos repository
}

func NewService(repos repository) *Service {
	return &Service{repos: repos}
}

func (s *Service) ReservationProducts(ctx context.Context, req models.ReservationProductsRequest) error {
	// сделать внутри одной транзакции, opeanapi, comments, postgis, testing
	productIDs, err := s.repos.ProductsIDsByPartNumbers(ctx, req.PartNumbers)
	if err != nil {
		return fmt.Errorf("error to get products: %v", err)
	}

	if len(productIDs) == 0 {
		return fmt.Errorf("products not found")
	}

	warehousesProducts, err := s.repos.WarehousesByProductIDs(ctx, productIDs, req.Latitude, req.Longitude)
	if err != nil {
		return fmt.Errorf("error to get warehouses: %v", err)
	}

	if len(warehousesProducts) == 0 {
		return fmt.Errorf("warehousesProducts not found")
	}

	reservation := make(map[int]int)
	for _, v := range warehousesProducts {
		if _, ok := reservation[v.ProductID]; !ok {
			reservation[v.ProductID] = v.WarehouseID
		}
	}

	err = s.repos.SetProductsToReserved(ctx, uuid.New(), reservation)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Products(ctx context.Context) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	products, err := s.repos.Products(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
