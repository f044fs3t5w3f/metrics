package handler

import (
	"net"
	"net/http"

	netTools "github.com/f044fs3t5w3f/metrics/pkg/net"

	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/f044fs3t5w3f/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

// GetRouter returns router for mertics services
// TODO: use service in every handler functions instead of storage
// Params:
// - storage: storage for metrics (to be removed)
// - service: service for metrics
// - key: key for sign middleware

func GetRouter(
	storage repository.Storage,
	service *service.Service,
	middlewares []func(http.Handler) http.Handler,
	allowedSubnet *net.IPNet,
) *chi.Mux {
	// service := service.NewService(storage)
	r := chi.NewRouter()
	if middlewares != nil {
		r.Use(middlewares...)
	}

	subnetMiddleware := netTools.GetCheckSubnetMiddleware(allowedSubnet)

	r.Get("/ping", ping(service))
	r.Post("/value/", GetJSON(storage))

	localSubnetRoutes := chi.NewRouter()
	localSubnetRoutes.Use(subnetMiddleware)
	localSubnetRoutes.Get("/value/{metricType}/{mericName}", Get(storage))
	localSubnetRoutes.Post("/update/", UpdateJSON(service))
	localSubnetRoutes.Post("/updates/", UpdatesJSON(service))
	localSubnetRoutes.Post("/update/{metricType}/{mericName}/{metricValue}", Update(service))
	r.Mount("/", localSubnetRoutes)

	r.Get("/", Index(storage))
	// TODO: use service everywhere
	return r
}
