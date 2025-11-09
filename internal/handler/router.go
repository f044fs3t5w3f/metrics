package handler

import (
	"database/sql"

	"github.com/f044fs3t5w3f/metrics/internal/compress"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

func GetRouter(storage repository.Storage, db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger, compress.Middleware)
	r.Get("/ping", ping(db))
	r.Post("/update/", UpdateJSON(storage))
	r.Post("/update/{metricType}/{mericName}/{metricValue}", Update(storage))
	r.Get("/value/{metricType}/{mericName}", Get(storage))
	r.Post("/value/", GetJSON(storage))
	r.Get("/", Index(storage))
	return r
}
