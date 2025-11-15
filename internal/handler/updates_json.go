package handler

import (
	"encoding/json"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func UpdatesJSON(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := []*models.Metrics{}
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)

		err := decoder.Decode(&metrics)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = storage.MultiUpdate(metrics)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
	}
}
