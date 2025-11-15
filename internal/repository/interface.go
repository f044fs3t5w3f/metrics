package repository

import "github.com/f044fs3t5w3f/metrics/internal/models"

type Storage interface {
	GetCounter(metricName string) (int64, error)
	GetGauge(metricName string) (float64, error)
	AddCounter(metricName string, value int64) error
	SetGauge(metricName string, value float64) error
	GetValuesList() ([]models.Metrics, error)
	MultiUpdate([]*models.Metrics) error
}
