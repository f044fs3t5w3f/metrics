package handler

import (
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

func GetRouter(storage repository.Storage) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	r.Post("/update/", UpdateJSON(storage))
	r.Post("/update/{metricType}/{mericName}/{metricValue}", Update(storage))
	r.Get("/value/{metricType}/{mericName}", Get(storage))
	r.Post("/value/", GetJSON(storage))
	return r
}
