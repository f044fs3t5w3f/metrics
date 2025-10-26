package main

import (
	"fmt"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func main() {
	parseFlags()
	parseEnv()
	storage := repository.NewMemStorage()
	addr := envRunAddr
	if addr == "" {
		addr = flagRunAddr
	}
	r := handler.GetRouter(storage)
	fmt.Println(addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		panic(err.Error())
	}
}
