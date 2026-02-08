package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/service"
)

// UpdateJSON updates metric with JSON request
// In case of incorrect request returns 400
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

		ctx := context.WithValue(r.Context(), service.CtxUserIP, r.RemoteAddr)
		err = s.UpdateMetric(ctx, metric)
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
