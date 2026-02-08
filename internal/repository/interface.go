// Repository represents the declaration of the Storage interface.

package repository

import "github.com/f044fs3t5w3f/metrics/internal/models"

// Storage defines an interface for storing and managing different kinds of metrics.
type Storage interface {
	// GetCounter gets the current value of a counter metric.
	GetCounter(metricName string) (int64, error)

	// GetGauge gets the current value of a gauge metric.
	GetGauge(metricName string) (float64, error)

	// AddCounter increases the value of a counter metric.
	AddCounter(metricName string, value int64) error

	// SetGauge sets a new value for a gauge metric.
	SetGauge(metricName string, value float64) error

	// GetValuesList lists all current metric values.
	GetValuesList() ([]models.Metrics, error)

	// MultiUpdate updates multiple metrics at once.
	MultiUpdate([]*models.Metrics) error

	// Ping checks if the storage is working properly.
	Ping() error
}
