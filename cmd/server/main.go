package main

import (
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func main() {
	storage := repository.NewMemStorage()
	r := handler.GetRouter(storage)
	err := http.ListenAndServe(`:8080`, r)
	if err != nil {
		panic(":(")
	}
}
