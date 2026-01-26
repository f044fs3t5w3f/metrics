package handler

import (
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/service"
	"github.com/go-chi/chi/v5"
)

func Update(service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type_ := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "mericName")
		merticValueStr := chi.URLParam(r, "metricValue")

		err := service.Update(type_, metricName, merticValueStr)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
	}
}
