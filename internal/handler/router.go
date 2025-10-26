package handler

import (
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

func GetRouter(storage repository.Storage) *chi.Mux {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Post("/update/{metricType}/{mericName}/{metricValue}", logger.RequestLogger(Update(storage)))
	r.Get("/value/{metricType}/{mericName}", logger.RequestLogger(Get(storage)))
	return r
}
