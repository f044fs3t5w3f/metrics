package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/crypto"
	"github.com/f044fs3t5w3f/metrics/internal/handler"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	dbRepo "github.com/f044fs3t5w3f/metrics/internal/repository/db"
	"github.com/f044fs3t5w3f/metrics/internal/repository/file"
	"github.com/f044fs3t5w3f/metrics/internal/service"
	"github.com/f044fs3t5w3f/metrics/internal/utils"
	"github.com/f044fs3t5w3f/metrics/pkg/compress"
	"github.com/f044fs3t5w3f/metrics/pkg/configuration"
	"github.com/f044fs3t5w3f/metrics/pkg/sign"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	// _ "net/http/pprof"
)

var retryPolicy []time.Duration = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	utils.PrintBuildInfo(os.Stdout, buildVersion, buildDate, buildCommit)
	config := Config{}
	err := configuration.ScanConfig(&config, nil)
	if err != nil {
		log.Fatalf("Config init: %s", err.Error())
	}
	err = logger.Initialize("INFO")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %s", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var storage repository.Storage

	if config.DatabaseDSN != "" {
		db, err := sql.Open("pgx", config.DatabaseDSN)
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
		storage = file.NewFileStorage(config.FileStorage, config.StoreInterval, config.RestoreOnBoot)
	}

	auditPublisher := audit.NewAuditPublisher(nil)
	if config.AuditFile != "" {
		auditPublisher.AddSubscriber(audit.NewRemoteAudit(config.AuditFile))
	}
	var fileAuditCleanup func()
	if config.AuditFile != "" {
		fileAudit, err := audit.NewFileAudit(ctx, config.AuditFile)
		fileAuditCleanup = fileAudit.Close
		if err == nil {
			auditPublisher.AddSubscriber(fileAudit)
		} else {
			logger.Log.Info("audit: cannot open file", zap.String("file", config.AuditFile))
		}
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGQUIT, syscall.SIGQUIT)

	service := service.NewService(storage, auditPublisher)
	service.AddCleanup(fileAuditCleanup)

	middlewares := []func(http.Handler) http.Handler{logger.RequestLogger}
	if config.CryptoFile != "" {
		privateKey, err := crypto.GetPrivateKey(config.CryptoFile)
		if err != nil {
			log.Fatalf("getPrivateKey: %s", err.Error())
		}
		if privateKey != nil {
			middlewares = append(middlewares, crypto.GetDecryptMiddleware(privateKey))
		}
	}
	if config.Key != "" {
		signMiddleware := sign.GetSignMiddleware(sign.GetSignFunc(config.Key))
		middlewares = append(middlewares, signMiddleware)
	}
	middlewares = append(middlewares, compress.Middleware, middleware.RealIP)

	var subnet *net.IPNet
	if config.TrustedSubnet != "" {
		_, subnet, err = net.ParseCIDR(config.TrustedSubnet)
		if err != nil {
			log.Fatalf("ParseCIDR: %s", err.Error())
		}
	}

	router := handler.GetRouter(storage, service, middlewares, subnet)

	srv := &http.Server{
		Addr:    config.RunAddr,
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	logger.Log.Info("Server has been started", zap.String("addr", config.RunAddr))
	go func() {
		logger.Log.Info("Starting http server", zap.String("addr", config.RunAddr))
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("couldn't start server", zap.Error(err))
		}
	}()

	sig := <-signals
	logger.Log.Info("shutting down", zap.String("signal", sig.String()))
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
	service.Close()
}
