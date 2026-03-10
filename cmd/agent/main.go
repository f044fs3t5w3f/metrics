package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	pb "github.com/f044fs3t5w3f/metrics/proto"
	"go.uber.org/zap"

	"github.com/f044fs3t5w3f/metrics/internal/agent"
	"github.com/f044fs3t5w3f/metrics/internal/crypto"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/utils"
	"github.com/f044fs3t5w3f/metrics/pkg/configuration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	buildVersion, buildDate, buildCommit string
)

type reportBatchFunc func(batch agent.MetricsBatch, wg *sync.WaitGroup)

func reportButchWithRPC(c pb.MetricsClient, batch agent.MetricsBatch) {
	metrics := make([]*pb.Metric, len(batch))
	for i, metric := range batch {
		metricBuilder := pb.Metric_builder{
			Id: metric.ID,
		}
		if metric.Delta != nil {
			metricBuilder.Delta = *metric.Delta
		}
		if metric.Value != nil {
			metricBuilder.Value = *metric.Value
		}
		switch metric.MType {
		case models.Counter:
			metricBuilder.Type = pb.Metric_COUNTER
		case models.Gauge:
			metricBuilder.Type = pb.Metric_GAUGE
		default:
			logger.Log.Error("Incorrect type for grpc", zap.String("type", metric.MType))
			return
		}
		metrics[i] = metricBuilder.Build()
	}
	request := pb.UpdateMetricsRequest_builder{
		Metrics: metrics,
	}.Build()
	_, err := c.UpdateMetrics(context.Background(), request)
	if err != nil {
		logger.Log.Error("grpc request failed", zap.Error(err))
	}
}

func getReportBatchFunc(config *Config) (reportBatchFunc, func(), error) {
	if config.RPCServer != "" {
		conn, err := grpc.NewClient(config.RPCServer,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithIdleTimeout(2*time.Second),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("grpc.NewClient, server %s: %w", config.RPCServer, err)
		}
		c := pb.NewMetricsClient(conn)
		closeFunc := func() {
			conn.Close()
		}
		return func(batch agent.MetricsBatch, wg *sync.WaitGroup) {
			wg.Add(1)
			reportButchWithRPC(c, batch)
			wg.Done()
		}, closeFunc, nil

	} else {
		var publicKey *rsa.PublicKey = nil
		if config.CryptoKeyPath != "" {
			var err error
			publicKey, err = crypto.GetPublicKey(config.CryptoKeyPath)
			if err != nil {
				return nil, nil, fmt.Errorf("crypto.GetPublicKey: %s", err.Error())
			}
		}
		return func(batch agent.MetricsBatch, wg *sync.WaitGroup) {
			agent.ReportBatch(config.Address, batch, config.Key, publicKey, wg)
		}, nil, nil
	}
}

func main() {
	utils.PrintBuildInfo(os.Stdout, buildVersion, buildDate, buildCommit)

	wg := &sync.WaitGroup{}
	canGoOn := atomic.Bool{}
	canGoOn.Store(true)

	config := Config{}
	err := configuration.ScanConfig(&config, nil)

	if err != nil {
		log.Fatalf("couldn't get config: %s", err.Error())
	}

	var pool chan struct{}
	if config.RateLimit > 0 {
		pool = make(chan struct{}, config.RateLimit)
	}

	lock := sync.Mutex{}
	var counter int64 = 0
	store := make([]agent.MetricsBatch, 0)
	err = logger.Initialize("INFO")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %s", err.Error())
	}

	reportBatch, cancel, err := getReportBatchFunc(&config)
	if cancel != nil {
		defer cancel()
	}
	if err != nil {
		log.Fatalf("getReportBatchFunc: %s", err.Error())
	}

	go func() {
		for {
			if !canGoOn.Load() {
				return
			}
			lock.Lock()
			if len(store) == 0 {
				lock.Unlock()
				continue
			}
			lastBatch := store[len(store)-1]
			lock.Unlock()
			go func() {
				if pool != nil {
					pool <- struct{}{}
					defer func() {
						<-pool
					}()
				}
				reportBatch(lastBatch, wg)
			}()

			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		}
	}()
	go func() {
		for {
			if !canGoOn.Load() {
				return
			}
			batch := agent.GetMetricsBatch(counter)
			counter += 1
			lock.Lock()
			store = append(store, batch)
			lock.Unlock()
			time.Sleep(time.Duration(config.PollInterval) * time.Second)
		}
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGQUIT, syscall.SIGQUIT)
	<-signals
	canGoOn.Store(false)

	wg.Wait()
}
