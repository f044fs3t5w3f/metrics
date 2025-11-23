package handler

import (
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func UpdateJSON(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metrics
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&metric)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case models.Gauge:
			if metric.Value == nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			err := storage.SetGauge(metric.ID, *metric.Value)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		case models.Counter:
			if metric.Delta == nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			err := storage.AddCounter(metric.ID, *metric.Delta)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
