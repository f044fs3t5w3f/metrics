package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/service"
)

// UpdatesJSON updates list of metric with JSON request
// In case of incorrect request returns 400
func UpdatesJSON(s *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := []*models.Metrics{}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), service.CtxUserIP, r.RemoteAddr)
		err = s.UpdateMetrics(ctx, metrics)

		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
