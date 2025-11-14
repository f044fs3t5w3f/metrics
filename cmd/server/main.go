package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	dbRepo "github.com/f044fs3t5w3f/metrics/internal/repository/db"
	"github.com/f044fs3t5w3f/metrics/internal/repository/file"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

func main() {
	parseFlags()
	parseEnv()
	err := logger.Initialize("INFO")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %s", err.Error())
	}

	fileStoragePath := envFileStoragePath
	if fileStoragePath == "" {
		fileStoragePath = flagFileStoragePath
	}

	var storeInterval int64
	if envStoreInterval != "" {
		var err error
		storeInterval, err = strconv.ParseInt(envStoreInterval, 10, 64)
		if err != nil {
			logger.Log.Fatal("couldn't parse store interval env", zap.Error(err))
		}
	} else {
		storeInterval = flagStoreInterval
	}

	var restore bool
	if envRestore != "" {
		restore, err = strconv.ParseBool(envRestore)
		if err != nil {
			logger.Log.Fatal("Incorrect value of environment variable RESTORE", zap.String("value", envRestore))
		}
	} else {
		restore = flagRestore
	}
	var storage repository.Storage

	var databaseParams string
	var db *sql.DB
	if envDatabaseParams != "" {
		databaseParams = envDatabaseParams
	} else {
		databaseParams = flagDatabaseParams
	}
	if databaseParams != "" {
		db, err = sql.Open("pgx", databaseParams)
		if err != nil {
			logger.Log.Fatal("couldn't connect to database", zap.Error(err))
		}

		err = migrateDB(db)

		if err != nil && err != migrate.ErrNoChange {
			logger.Log.Fatal("couldn't migrate database", zap.Error(err))
		}

		storage = dbRepo.NewDBStorage(db)
	}
	if storage == nil {
		storage = file.NewFileStorage(fileStoragePath, storeInterval, restore)

	}
	addr := envRunAddr
	if addr == "" {
		addr = flagRunAddr
	}

	r := handler.GetRouter(storage, db)
	logger.Log.Info("Server has been started", zap.String("addr", addr))
	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Log.Fatal("couldn't start server", zap.Error(err))
	}
}
