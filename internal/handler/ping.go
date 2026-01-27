package handler

import (
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/service"
)

func ping(service *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := service.Ping()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
