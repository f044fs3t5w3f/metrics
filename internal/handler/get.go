package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

// Get returns metric by type and name. Returns plain text
func Get(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type_ := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "mericName")

		switch type_ {
		case models.Gauge:
			value, err := storage.GetGauge(metricName)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			io.WriteString(w, fmt.Sprintf("%g", value))

		case models.Counter:
			value, err := storage.GetCounter(metricName)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			io.WriteString(w, fmt.Sprintf("%d", value))
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
