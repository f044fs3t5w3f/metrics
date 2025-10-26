package repository

import (
	"errors"
	"sync"
)

var ErrNoValue = errors.New("THERE IS NO VALUE")

func NewMemStorage() Storage {
	return &memStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

type memStorage struct {
	lock    sync.RWMutex
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *memStorage) AddCounter(metricName string, value int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Counter[metricName] += value
}

func (m *memStorage) SetGauge(metricName string, value float64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Gauge[metricName] = value
}

func (m *memStorage) GetCounter(metricName string) (int64, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok := m.Counter[metricName]
	if !ok {
		return 0, ErrNoValue
	}
	return value, nil
}

func (m *memStorage) GetGauge(metricName string) (float64, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	value, ok := m.Gauge[metricName]
	if !ok {
		return 0, ErrNoValue
	}
	return value, nil
}
