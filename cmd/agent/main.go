package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/agent"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
)

var (
	buildVersion, buildDate, buildCommit string
)

func printBuildInfo(w io.Writer) {
	v := func(val string) string {
		if val != "" {
			return val
		}
		return "N/A"
	}
	fmt.Fprintf(w, "Build version: %s\n", v(buildVersion))
	fmt.Fprintf(w, "Build date: %s\n", v(buildDate))
	fmt.Fprintf(w, "Build commit: %s\n", v(buildCommit))
}

func main() {
	printBuildInfo(os.Stdout)
	parseFlags()
	parseEnv()

	pollInterval := envPollInterval
	if pollInterval == 0 {
		pollInterval = flagPollInterval
	}

	key := envKey
	if key == "" {
		key = flagKey
	}

	reportInterval := envReportInterval
	if reportInterval == 0 {
		reportInterval = flagReportInterval
	}

	addr := envRunAddr
	if addr == "" {
		addr = flagEndpointAddr
	}

	rateLimit := envRateLimit
	if envRateLimit == 0 {
		rateLimit = flagRateLimit
	}
	var pool chan struct{}
	if rateLimit > 0 {
		pool = make(chan struct{}, rateLimit)
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
			go func() {
				if pool != nil {
					pool <- struct{}{}
					defer func() {
						<-pool
					}()
				}
				agent.ReportBatch(addr, lastBatch, key)
			}()

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
