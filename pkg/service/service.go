package service

import (
	"context"
	"math"
	"time"

	"github.com/Hymiside/lamoda-api/pkg/models"
)

type warehouseDistance struct {
	ID   int
	Dist float64
}

type repository interface {
	Products(ctx context.Context) ([]models.Product, error)
	ProductsIDsByPartNumbers(ctx context.Context, partNumbers []models.ReservationItems) ([]int, error)
	WarehouseIDsByProductID(ctx context.Context, productIDs []int) ([]int, error)
	WarehousesByIDs(ctx context.Context, warehouseIDs []int) ([]models.Warehouse, error)
	SetProductsToReserved(ctx context.Context, warehouseID int, productIDs []int) error
}

type Service struct {
	repos repository
}

func NewService(repos repository) *Service {
	return &Service{repos: repos}
}

func (s *Service) distanceToWarehouse(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371

	phi1, phi2 := lat1*math.Pi/180, lat2*math.Pi/180
	dPhi, dLambda := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180

	a := math.Sin(dPhi/2)*math.Sin(dPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(dLambda/2)*math.Sin(dLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}

func (s *Service) ReservationProducts(ctx context.Context, req models.ReservationProductsRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	productIDs, err := s.repos.ProductsIDsByPartNumbers(ctx, req.Items)
	if err != nil {
		return err
	}

	warehouseIDs, err := s.repos.WarehouseIDsByProductID(ctx, productIDs)
	if err != nil {
		return err
	}
	warehouses, err := s.repos.WarehousesByIDs(ctx, warehouseIDs)
	if err != nil {
		return err
	}

	ch := make(chan warehouseDistance, len(warehouses))
	for _, wh := range warehouses {
		go func(warehouse models.Warehouse) {
			ch <- warehouseDistance{
				ID:   warehouse.ID,
				Dist: s.distanceToWarehouse(warehouse.Latitude, warehouse.Longitude, req.Latitude, req.Longitude),
			}
		}(wh)
	}

	var whID int
	minDistanceToWarehouse := math.MaxFloat64
	for i := 0; i < len(warehouses); i++ {
		wh := <-ch
		if wh.Dist < minDistanceToWarehouse {
			minDistanceToWarehouse = wh.Dist
			whID = wh.ID
		}
	}

	err = s.repos.SetProductsToReserved(ctx, whID, productIDs)
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
