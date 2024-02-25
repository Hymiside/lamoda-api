package service

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/google/uuid"
)

type repository interface {
	Products(ctx context.Context) ([]models.Product, error)
	ProductsIDsByPartNumbers(ctx context.Context, partNumbers []string) ([]int, error)
	WarehousesByProductIDs(ctx context.Context, productIDs []int) ([]models.WarehouseProductID, error)
	SetProductsToReserved(ctx context.Context, reservationID uuid.UUID, warehousesProducts []models.ReservationProducts) error
}

type reservationData struct {
	mx    sync.RWMutex
	wp    map[int]models.WarehouseInfo
	avail map[int]int
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

	productIDs, err := s.repos.ProductsIDsByPartNumbers(ctx, req.PartNumbers)
	if err != nil {
		return err
	}

	warehousesProducts, err := s.repos.WarehousesByProductIDs(ctx, productIDs)
	if err != nil {
		return fmt.Errorf("error to get warehouses: %v", err)
	}

	rd := reservationData{
		wp:    make(map[int]models.WarehouseInfo),
		avail: make(map[int]int),
	}

	wg := sync.WaitGroup{}
	wg.Add(len(warehousesProducts))

	for _, wh := range warehousesProducts {
		go func(warehouse models.WarehouseProductID) {
			rd.mx.Lock()
			defer func() {
				rd.mx.Unlock()
				wg.Done()
			}()

			if _, ok := rd.avail[warehouse.Warehouse.ID]; !ok {
				rd.avail[warehouse.Warehouse.ID] = warehouse.Quantity
			}

			distance := s.distanceToWarehouse(req.Latitude, req.Longitude, warehouse.Warehouse.Latitude, warehouse.Warehouse.Longitude)

			if _, ok := rd.wp[warehouse.ProductID]; !ok && rd.avail[warehouse.Warehouse.ID] > 0 {
				rd.wp[warehouse.ProductID] = models.WarehouseInfo{
					WarehouseID: warehouse.Warehouse.ID,
					Dist:        distance,
				}
				rd.avail[warehouse.Warehouse.ID]--
			} else {
				if rd.wp[warehouse.ProductID].Dist > distance && rd.avail[warehouse.Warehouse.ID] > 0 {
					rd.wp[warehouse.ProductID] = models.WarehouseInfo{
						WarehouseID: warehouse.Warehouse.ID,
						Dist:        distance,
					}
					rd.avail[warehouse.Warehouse.ID]--
				}
			}
		}(wh)
	}

	wg.Wait()

	var reservation []models.ReservationProducts
	for i, v := range rd.wp {
		reservation = append(
			reservation,
			models.ReservationProducts{
				ProductID:   i,
				WarehouseID: v.WarehouseID,
			})
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
