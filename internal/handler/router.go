package handler

import (
	"github.com/f044fs3t5w3f/metrics/internal/compress"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/f044fs3t5w3f/metrics/internal/sign"
	"github.com/go-chi/chi/v5"
)

func GetRouter(storage repository.Storage, key string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	if key != "" {
		signMiddleware := sign.GetSignMiddleware(sign.GetSignFunc(key))
		r.Use(signMiddleware)
	}
	r.Use(compress.Middleware)
	r.Get("/ping", ping(storage))
	r.Post("/update/", UpdateJSON(storage))
	r.Post("/updates/", UpdatesJSON(storage))
	r.Post("/update/{metricType}/{mericName}/{metricValue}", Update(storage))
	r.Get("/value/{metricType}/{mericName}", Get(storage))
	r.Post("/value/", GetJSON(storage))
	r.Get("/", Index(storage))
	return r
}
