package agent

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type MetricsBatch []*models.Metrics

func memStatsToMap(memStats *runtime.MemStats) map[string]float64 {
	bytes, _ := json.Marshal(memStats)
	metrics := make(map[string]float64)
	json.Unmarshal(bytes, &metrics)
	return metrics
}

func getMainMetrics() (map[string]float64, error) {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)
	return memStatsToMap(memStats), nil
}

func getAdditionalMetrics() (map[string]float64, error) {
	result := make(map[string]float64)

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	result["TotalMemory"] = float64(vm.Total)
	result["FreeMemory"] = float64(vm.Free)

	cpuPercents, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	for i, v := range cpuPercents {
		result[fmt.Sprintf("CPUutilization%d", i+1)] = v
	}

	return result, nil
}

func GetMetricsBatch(counter int64) MetricsBatch {
	metricsFuncs := [](func() (map[string]float64, error)){getMainMetrics, getAdditionalMetrics}
	var wg sync.WaitGroup
	var lock sync.Mutex
	metrics := make(map[string]float64)
	for _, f := range metricsFuncs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			newMetrics, err := f()
			if err != nil {
				return
			}
			lock.Lock()
			defer lock.Unlock()
			for k, v := range newMetrics {
				metrics[k] = v
			}
		}()
	}
	wg.Wait()
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
		ID:    "RandomValue",
		MType: models.Gauge,
		Value: &ramdomValue,
	})
	return batch
}
