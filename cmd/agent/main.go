package main

import (
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/agent"
)

func main() {
	parseFlags()
	pollInterval := time.Duration(flagPollInterval) * time.Second
	reportInterval := time.Duration(flagReportInterval) * time.Second
	lock := sync.Mutex{}
	var counter int64 = 0
	store := make([]agent.MetricsBatch, 0)
	go func() {
		for {
			lock.Lock()
			if len(store) == 0 {
				lock.Unlock()
				continue
			}
			lock.Unlock()
			lastBatch := store[len(store)-1]
			agent.ReportBatch(flagEndpointAddr, lastBatch)
			time.Sleep(reportInterval)
		}
	}()
	for {
		batch := agent.GetMetricsBatch(counter)
		counter += 1
		lock.Lock()
		store = append(store, batch)
		lock.Unlock()
		time.Sleep(pollInterval)
	}
}
