package file

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
	"github.com/f044fs3t5w3f/metrics/internal/repository"
	"github.com/f044fs3t5w3f/metrics/internal/repository/memory"
	"go.uber.org/zap"
)

type fileStorage struct {
	repository.Storage
	saveToFileOnChange bool
	fileStoragePath    string
	lock               sync.Mutex
}

func NewFileStorage(fileStoragePath string, storeInterval int64, restore bool) repository.Storage {
	memStorage := memory.NewMemStorage()
	storage := &fileStorage{
		Storage:            memStorage,
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

func (fs *fileStorage) restoreFromFile(r io.Reader) error {
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
			fs.Storage.AddCounter(metricName, *metric.Delta)
		case models.Gauge:
			fs.Storage.SetGauge(metricName, *metric.Value)
		}
	}
	return nil
}

func (fs *fileStorage) saveToFile() {
	logger.Log.Info("saving file...", zap.String("file", fs.fileStoragePath))

	fs.lock.Lock()
	defer fs.lock.Unlock()
	file, err := os.OpenFile(fs.fileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Log.Warn("error while opening file", zap.String("file", fs.fileStoragePath), zap.Error(err))
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	metrics := fs.GetValuesList()
	err = encoder.Encode(metrics)
	if err != nil {
		logger.Log.Warn("error while encoding metrics", zap.String("file", fs.fileStoragePath), zap.Error(err))
	} else {
		logger.Log.Info("file was saved", zap.String("file", fs.fileStoragePath))
	}
}

func (fs *fileStorage) AddCounter(metricName string, value int64) {
	fs.Storage.AddCounter(metricName, value)
	if fs.saveToFileOnChange {
		go fs.saveToFile()
	}
}

func (fs *fileStorage) SetGauge(metricName string, value float64) {
	fs.Storage.SetGauge(metricName, value)
	if fs.saveToFileOnChange {
		go fs.saveToFile()
	}
}
