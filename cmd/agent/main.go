package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

const pollInterval = 2 * time.Second
const reportInterval = 10 * time.Second

func memStatsToMap(memStats *runtime.MemStats) map[string]float64 {
	bytes, _ := json.Marshal(memStats)
	metrics := make(map[string]float64)
	json.Unmarshal(bytes, &metrics)
	return metrics
}

type metricsBatch []*models.Metrics

func getMetricsBatch(counter int64) metricsBatch {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	metrics := memStatsToMap(memStats)
	batch := make(metricsBatch, 0)
	for metricName, metricValue := range metrics {
		metric := models.Metrics{
			ID:    metricName,
			MType: models.Gauge,
			Value: &metricValue,
		}
		batch = append(batch, &metric)
	}
	batch = append(batch, &models.Metrics{
		ID:    "PollCount",
		MType: models.Counter,
		Delta: &counter,
	})
	ramdomValue := rand.Float64()
	batch = append(batch, &models.Metrics{
		ID:    "PollCount",
		MType: models.Gauge,
		Value: &ramdomValue,
	})
	counter += 1
	return batch
}

func report(batch metricsBatch) {
	baseUrl := "http://localhost:8080/update/"
	for _, metric := range batch {
		var url string
		switch metric.MType {
		case models.Counter:
			url = fmt.Sprintf("%scounter/%s/%d", baseUrl, metric.ID, *metric.Delta)
		case models.Gauge:
			url = fmt.Sprintf("%sgauge/%s/%f", baseUrl, metric.ID, *metric.Value)
		default:
			continue
		}
		fmt.Println(url)
		response, err := http.Post(url, "", nil)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(response.Status)
	}
}

func main() {
	lock := sync.Mutex{}
	var counter int64 = 0
	store := make([]metricsBatch, 0)
	go func() {
		for {
			lock.Lock()
			if len(store) == 0 {
				lock.Unlock()
				continue
			}
			lock.Unlock()
			lastBatch := store[len(store)-1]
			report(lastBatch)
			time.Sleep(reportInterval)
		}
	}()
	for {
		batch := getMetricsBatch(counter)
		lock.Lock()
		store = append(store, batch)
		lock.Unlock()
		time.Sleep(pollInterval)
	}
}
