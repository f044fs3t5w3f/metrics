package handler

import (
	"net/http"
	"strconv"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/go-chi/chi/v5"
)

func Update(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type_ := chi.URLParam(r, "metricType")
		metricName := chi.URLParam(r, "mericName")
		merticValueStr := chi.URLParam(r, "metricValue")

		switch type_ {
		case models.Gauge:
			merticParsed, err := strconv.ParseFloat(merticValueStr, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			storage.SetGauge(metricName, merticParsed)
		case models.Counter:
			merticParsed, err := strconv.ParseInt(merticValueStr, 0, 64)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			storage.AddCounter(metricName, merticParsed)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
