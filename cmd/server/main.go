package main

import (
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/repository"

	"github.com/go-chi/chi/v5"
)

func main() {

	memStorage := repository.NewMemStorage()
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Post("/update/{metricType}/{mericName}/{metricValue}", handler.Update(memStorage))
	r.Get("/value/{metricType}/{mericName}/", handler.Get(memStorage))
	// r.Use(middleware.Recoverer)
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(":(")
	}
}
