package repository

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
	"go.uber.org/zap"
)

var ErrNoValue = errors.New("THERE IS NO VALUE")

func NewMemStorage(fileStoragePath string, storeInterval int64, restore bool) Storage {
	storage := &memStorage{
		gauge:              make(map[string]float64),
		counter:            make(map[string]int64),
		saveToFileOnChange: storeInterval == 0,
		fileStoragePath:    fileStoragePath,
	}
	if restore {
		file, err := os.Open(fileStoragePath)
		if err != nil {
			logger.Log.Warn("failed to open file to restore settings", zap.String("file", fileStoragePath))
		}
		defer file.Close()
		storage.restoreFromFile(file)
	}
	if storeInterval != 0 {
		go func() {
			for {
				time.Sleep(time.Duration(storeInterval) * time.Second)
				storage.saveToFile()
			}
		}()
	}
	return storage
}
func NewMemStorageWithoutFile() Storage {
	storage := &memStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
	return storage
}

type memStorage struct {
	lock               sync.RWMutex
	gauge              map[string]float64
	counter            map[string]int64
	saveToFileOnChange bool
	fileStoragePath    string
}

func (m *memStorage) restoreFromFile(r io.Reader) error {
	metrics := make([]models.Metrics, 0)
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&metrics)
	if err != nil {
		return err
	}

	for _, metric := range metrics {
		metricName := metric.ID
		switch metric.MType {
		case models.Counter:
			m.counter[metricName] = *metric.Delta
		case models.Gauge:
			m.gauge[metricName] = *metric.Value
		}
	}
	return nil
}

func (m *memStorage) saveToFile() {
	logger.Log.Info("saving file...", zap.String("file", m.fileStoragePath))

	m.lock.RLock()
	defer m.lock.RUnlock()
	file, err := os.OpenFile(m.fileStoragePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Log.Warn("error while opening file", zap.String("file", m.fileStoragePath), zap.Error(err))
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	metrics := m.GetValuesList()
	err = encoder.Encode(metrics)
	if err != nil {
		logger.Log.Warn("error while encoding metrics", zap.String("file", m.fileStoragePath), zap.Error(err))
	} else {
		logger.Log.Info("file was saved", zap.String("file", m.fileStoragePath))
	}
}

// GetValuesList implements Storage.
func (m *memStorage) GetValuesList() []models.Metrics {
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
	return list
}

func (m *memStorage) AddCounter(metricName string, value int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.counter[metricName] += value
	if m.saveToFileOnChange {
		go m.saveToFile()
	}
}

func (m *memStorage) SetGauge(metricName string, value float64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.gauge[metricName] = value
	if m.saveToFileOnChange {
		go m.saveToFile()
	}
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
