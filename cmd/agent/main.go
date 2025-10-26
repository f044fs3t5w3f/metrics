package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/agent"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
)

func main() {
	parseFlags()
	parseEnv()

	pollInterval := envPollInterval
	if pollInterval == 0 {
		pollInterval = flagPollInterval
	}

	reportInterval := envReportInterval
	if reportInterval == 0 {
		reportInterval = flagReportInterval
	}

	addr := envRunAddr
	if addr == "" {
		addr = flagEndpointAddr
	}

	lock := sync.Mutex{}
	var counter int64 = 0
	store := make([]agent.MetricsBatch, 0)
	err := logger.Initialize("INFO")
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
			fmt.Println("send")
			agent.ReportBatch(addr, lastBatch)
			time.Sleep(time.Duration(reportInterval) * time.Second)
		}
	}()
	for {
		batch := agent.GetMetricsBatch(counter)
		counter += 1
		lock.Lock()
		store = append(store, batch)
		lock.Unlock()
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
