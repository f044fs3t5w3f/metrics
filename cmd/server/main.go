package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	dbRepo "github.com/f044fs3t5w3f/metrics/internal/repository/db"
	"github.com/f044fs3t5w3f/metrics/internal/repository/file"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var retryPolicy []time.Duration = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

func main() {
	config, err := getConfig()
	if err != nil {
		log.Fatalf("Config init: %s", err.Error())
	}
	err = logger.Initialize("INFO")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %s", err.Error())
	}

	var storage repository.Storage

	if config.databaseParams != "" {
		db, err := sql.Open("pgx", config.databaseParams)
		if err != nil {
			logger.Log.Fatal("couldn't connect to database", zap.Error(err))
		}

		err = migrateDB(db)

		if err != nil && err != migrate.ErrNoChange {
			logger.Log.Fatal("couldn't migrate database", zap.Error(err))
		}

		storage = dbRepo.NewDBStorage(db, retryPolicy)
	}
	if storage == nil {
		storage = file.NewFileStorage(config.fileStoragePath, config.storeInterval, config.restore)
	}

	r := handler.GetRouter(storage)
	logger.Log.Info("Server has been started", zap.String("addr", config.runAddr))
	err = http.ListenAndServe(config.runAddr, r)
	if err != nil {
		logger.Log.Fatal("couldn't start server", zap.Error(err))
	}
}
