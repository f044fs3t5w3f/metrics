package handler

import (
	"github.com/f044fs3t5w3f/metrics/internal/compress"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/f044fs3t5w3f/metrics/internal/service"
	"github.com/f044fs3t5w3f/metrics/internal/sign"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func GetRouter(storage repository.Storage, service *service.Service, key string) *chi.Mux {
	// service := service.NewService(storage)
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)
	if key != "" {
		signMiddleware := sign.GetSignMiddleware(sign.GetSignFunc(key))
		r.Use(signMiddleware)
	}
	r.Use(compress.Middleware)
	r.Use(middleware.RealIP)
	r.Get("/ping", ping(service))
	r.Post("/update/", UpdateJSON(service))
	r.Post("/updates/", UpdatesJSON(service))
	r.Post("/update/{metricType}/{mericName}/{metricValue}", Update(service))
	r.Get("/value/{metricType}/{mericName}", Get(storage))
	r.Post("/value/", GetJSON(storage))
	r.Get("/", Index(storage))
	// TODO: use service everywhere
	return r
}
