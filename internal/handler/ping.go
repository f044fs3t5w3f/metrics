package handler

import (
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func ping(storage repository.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := storage.Ping()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
