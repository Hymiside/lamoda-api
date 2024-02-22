package service

import (
	"context"
	"sync"
	"math"
	"time"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/Hymiside/lamoda-api/pkg/repository"
)

type storeService struct {
	repos *repository.Repository
}

func newStoreService(repos *repository.Repository) *storeService {
	return &storeService{repos: repos}
}

func (s *storeService) ReservationProducts(_ context.Context, data models.ProductReservationRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	productIDs, err := s.repos.S.ProductsIDsByPartNumbers(ctx, data.PartNumbers)
	if err != nil {
		return err
	}
	warehouseIDs, err := s.repos.S.WarehouseIDsByProductID(ctx, productIDs)
	if err != nil {
		return err
	}
	warehouses, err := s.repos.S.WarehousesByIDs(ctx, warehouseIDs)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(warehouses))

	ch := make(chan [2]interface{}, len(warehouses))
	for _, warehouse := range warehouses {

		go func(wh models.Warehouse) {
			defer wg.Done()
			ch <- [2]interface{}{wh.ID, s.distanceToWarehouse(wh.Lat, wh.Long, data.Lat, data.Long)}
		}(warehouse)

	}
	wg.Wait()

	minDistanceToWarehouse := [2]interface{}{0, math.MaxFloat64}
	for i := 0; i < len(warehouses); i++ {
		res := <-ch
		if res[1].(float64) < minDistanceToWarehouse[1].(float64) {
			minDistanceToWarehouse = res
		}
	}

	err = s.repos.S.SetProductsToReserved(ctx, minDistanceToWarehouse[0].(int), productIDs)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *storeService) distanceToWarehouse(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371

	phi1, phi2 := lat1*math.Pi/180, lat2*math.Pi/180
	dPhi, dLambda := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180

	a := math.Sin(dPhi/2)*math.Sin(dPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(dLambda/2)*math.Sin(dLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c

	return distance
}
