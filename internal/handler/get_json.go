package handler

import (
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func GetJson(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric *models.Metrics
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&metric)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case models.Gauge:
			value, err := storage.GetGauge(metric.ID)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			metric.Value = &value
		case models.Counter:
			value, err := storage.GetCounter(metric.ID)
			if err != nil {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			metric.Delta = &value
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		encoder := json.NewEncoder(w)
		err = encoder.Encode(metric)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}
