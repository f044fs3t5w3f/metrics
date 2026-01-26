package service

import "github.com/f044fs3t5w3f/metrics/internal/repository"

type Service struct {
	storage repository.Storage
}

func NewService(storage repository.Storage) *Service {
	return &Service{
		storage: storage,
	}
}
