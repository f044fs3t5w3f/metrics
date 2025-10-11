package repository

import (
	"errors"
)

type Storage interface {
	GetCounter(metricName string) (int64, error)
	GetGauge(metricName string) (float64, error)
	AddCounter(metricName string, value int64)
	SetGauge(metricName string, value float64)
}

func NewMemStorage() Storage {
	return &memStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

var ErrNoValue = errors.New("THERE IS NO VALUE")

type memStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func (m *memStorage) AddCounter(metricName string, value int64) {
	m.Counter[metricName] += value

}

func (m *memStorage) GetCounter(metricName string) (int64, error) {
	value, ok := m.Counter[metricName]
	if !ok {
		return 0, ErrNoValue
	}
	return value, nil
}

func (m *memStorage) GetGauge(metricName string) (float64, error) {
	value, ok := m.Gauge[metricName]
	if !ok {
		return 0, ErrNoValue
	}
	return value, nil
}

func (m *memStorage) SetGauge(metricName string, value float64) {
	m.Gauge[metricName] = value
}
