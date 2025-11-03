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
		// data, _ := io.ReadAll(r.Body)
		// dataS := string(data)
		// fmt.Println(dataS)
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
			storage.SetGauge(metric.ID, *metric.Value)
			return
		case models.Counter:
			if metric.Delta == nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			storage.AddCounter(metric.ID, *metric.Delta)
			return
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
