package agent

import (
	"encoding/json"
	"math/rand/v2"
	"runtime"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

type MetricsBatch []*models.Metrics

func memStatsToMap(memStats *runtime.MemStats) map[string]float64 {
	bytes, _ := json.Marshal(memStats)
	metrics := make(map[string]float64)
	json.Unmarshal(bytes, &metrics)
	return metrics
}

func GetMetricsBatch(counter int64) MetricsBatch {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	metrics := memStatsToMap(memStats)
	batch := make(MetricsBatch, 0)
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
	return batch
}
