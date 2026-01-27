package service

import (
	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

type ctxKey string

const CtxUserIP ctxKey = "userIP"

type Service struct {
	storage repository.Storage
	audit   audit.Audit
}

func NewService(storage repository.Storage, audit audit.Audit) *Service {
	return &Service{
		storage: storage,
		audit:   audit,
	}
}
