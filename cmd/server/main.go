package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

func main() {
	parseFlags()
	parseEnv()

	fileStoragePath := envFileStoragePath
	if fileStoragePath == "" {
		fileStoragePath = flagFileStoragePath
	}

	var storeInterval int64
	if envStoreInterval != "" {
		var err error
		storeInterval, err = strconv.ParseInt(envStoreInterval, 10, 64)
		if err != nil {
			log.Fatalf("couldn't parse store interval env: %s", err.Error())
		}
	} else {
		storeInterval = flagStoreInterval
	}

	var restore bool
	if envRestore != "" {
		switch strings.ToLower(envRestore) {
		case "true":
			restore = true
		case "false":
			restore = false
		default:
			log.Fatalf("Incorrect value of environment variable RESTORE: %s. Only true/false are avaliable", envRestore)
		}
	} else {
		restore = flagRestore
	}

	storage := repository.NewMemStorage(fileStoragePath, storeInterval, restore)
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
