package service

import (
	"context"

	"github.com/Hymiside/lamoda-api/pkg/models"
	"github.com/Hymiside/lamoda-api/pkg/repository"
)

type store interface {
	ReservationProducts(_ context.Context, data models.ProductReservationRequest) error
}

type Service struct {
	S store
}

func NewService(repos *repository.Repository) *Service {
	return &Service{S: newStoreService(repos)}
}