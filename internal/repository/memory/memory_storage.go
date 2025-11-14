package memory

import (
	"errors"
	"sync"

	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
)

var ErrNoValue = errors.New("THERE IS NO VALUE")

func NewMemStorage() repository.Storage {
	storage := &memStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
	return storage
}

type memStorage struct {
	lock    sync.RWMutex
	gauge   map[string]float64
	counter map[string]int64
}

func (m *memStorage) GetValuesList() ([]models.Metrics, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	list := make([]models.Metrics, 0, len(m.gauge)+len(m.counter))
	for metricName, metricValue := range m.gauge {
		metric := models.Metrics{
			MType: models.Gauge,
			ID:    metricName,
			Value: &metricValue,
		}
		list = append(list, metric)
	}
	for metricName, metricValue := range m.counter {
		metric := models.Metrics{
			MType: models.Counter,
			ID:    metricName,
			Delta: &metricValue,
		}
		list = append(list, metric)
	}
	return list, nil
}

func (m *memStorage) AddCounter(metricName string, value int64) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.counter[metricName] += value
	return nil
}

func (m *memStorage) SetGauge(metricName string, value float64) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.gauge[metricName] = value
	return nil
}

func (m *memStorage) GetCounter(metricName string) (int64, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok := m.counter[metricName]
	if !ok {
		return 0, ErrNoValue
	}
	return value, nil
}

func (m *memStorage) GetGauge(metricName string) (float64, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok := m.gauge[metricName]
	if !ok {
		return 0, ErrNoValue
	}
	return value, nil
}
