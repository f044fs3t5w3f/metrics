package main

import (
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func main() {

	memStorage := models.MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handler.Update(memStorage))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(":(")
	}
}
