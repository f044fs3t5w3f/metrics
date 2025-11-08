package agent

import (
	"errors"
	"fmt"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

var ErrEmptyValue = errors.New("EMPTY VALUE")
var ErrEmptyDelta = errors.New("EMPTY DELTA")

// Deprecated
func getURL(host string, metric *models.Metrics) (string, error) {
	baseURL := fmt.Sprintf("http://%s/update/", host)
	switch metric.MType {
	case models.Counter:
		if metric.Delta == nil {
			return "", ErrEmptyDelta
		}
		return fmt.Sprintf("%scounter/%s/%d", baseURL, metric.ID, *metric.Delta), nil
	case models.Gauge:
		if metric.Value == nil {
			return "", ErrEmptyValue
		}
		return fmt.Sprintf("%sgauge/%s/%g", baseURL, metric.ID, *metric.Value), nil
	default:
		return "", errors.New("incorrect type")
	}
}
