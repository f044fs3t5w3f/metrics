package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
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
	err := logger.Initialize("INFO")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %s", err.Error())
	}
	fmt.Println(addr)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("couldn't start server: %s", err)
	}
}
