package main

import (
	"crypto/rsa"
	"log"
	"os"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/agent"
	cfg "github.com/f044fs3t5w3f/metrics/internal/agent/config"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/utils"
)

var (
	buildVersion, buildDate, buildCommit string
)

func main() {
	utils.PrintBuildInfo(os.Stdout, buildVersion, buildDate, buildCommit)

	config, err := cfg.GetConfig()
	if err != nil {
		log.Fatalf("couldn't get config: %s", err.Error())
	}

	var pool chan struct{}
	if config.RateLimit > 0 {
		pool = make(chan struct{}, config.RateLimit)
	}

	var publicKey *rsa.PublicKey = nil
	if config.CryptoKeyPath != "" {

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
				agent.ReportBatch(config.RunAddr, lastBatch, config.Key, publicKey)
			}()

			time.Sleep(time.Duration(config.ReportInterval) * time.Second)
		}
	}()
	for {
		batch := agent.GetMetricsBatch(counter)
		counter += 1
		lock.Lock()
		store = append(store, batch)
		lock.Unlock()
		time.Sleep(time.Duration(config.PollInterval) * time.Second)
	}
}
