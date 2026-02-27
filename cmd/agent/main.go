package main

import (
	"crypto/rsa"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/agent"
	"github.com/f044fs3t5w3f/metrics/internal/crypto"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/utils"
	"github.com/f044fs3t5w3f/metrics/pkg/configuration"
)

var (
	buildVersion, buildDate, buildCommit string
)

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

	var publicKey *rsa.PublicKey = nil
	if config.CryptoKeyPath != "" {
		publicKey, err = crypto.GetPublicKey(config.CryptoKeyPath)
		if err != nil {
			log.Fatalf("getPublicKey: %s", err.Error())
		}
	}

	lock := sync.Mutex{}
	var counter int64 = 0
	store := make([]agent.MetricsBatch, 0)
	err = logger.Initialize("INFO")
	if err != nil {
		log.Fatalf("couldn't initialize logger: %s", err.Error())
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
				agent.ReportBatch(config.Address, lastBatch, config.Key, publicKey, wg)
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
