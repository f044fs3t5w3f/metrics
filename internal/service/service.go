package service

import (
	"sync"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

type ctxKey string

const CtxUserIP ctxKey = "userIP"

type Service struct {
	storage  repository.Storage
	audit    audit.Audit
	cleanups []func()
}

func NewService(storage repository.Storage, audit audit.Audit) *Service {
	return &Service{
		storage: storage,
		audit:   audit,
	}
}

func (s *Service) AddCleanup(cleanup func()) {
	if cleanup != nil {
		s.cleanups = append(s.cleanups, cleanup)
	}
}

// Close method executes all the registered cleanup functions.
// It waits until all the cleanup functions are completed.
func (s *Service) Close() {
	wg := sync.WaitGroup{}
	for _, cleanup := range s.cleanups {
		wg.Add(1)
		go func() {
			cleanup()
			wg.Done()
		}()
	}
	wg.Wait()
	logger.Log.Info("All the cleanups were completed")
}
