package main

import (
	"net/http"
)

func main() {

	memStorage := memStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", handleUpdate(memStorage))
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(":(")
	}
}
