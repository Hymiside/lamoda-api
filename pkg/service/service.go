package service

import "github.com/Hymiside/lamoda-api/pkg/repository"

type store interface {}

type Service struct {
	s store
}

func NewService(repos *repository.Repository) *Service {
	return &Service{s: newStoreService(repos)}
}