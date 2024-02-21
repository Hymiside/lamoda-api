package service

import "github.com/Hymiside/lamoda-api/pkg/repository"

type storeService struct {
	repos *repository.Repository
}

func newStoreService(repos *repository.Repository) *storeService {
	return &storeService{repos: repos}
}