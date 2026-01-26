package handler

import (
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/service"
)

func UpdateJSON(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric models.Metrics
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&metric)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = s.UpdateMetric(metric)
		if err != nil {
			switch err {
			case service.ErrBadValue:
				http.Error(w, "Bad request", http.StatusBadRequest)
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
}
