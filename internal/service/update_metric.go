package service

import (
	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func (s *Service) UpdateMetric(metric models.Metrics) error {

	switch metric.MType {
	case models.Gauge:
		if metric.Value == nil {
			return ErrBadValue
		}
		err := s.storage.SetGauge(metric.ID, *metric.Value)
		return err
	case models.Counter:
		if metric.Delta == nil {
			return ErrBadValue
		}
		err := s.storage.AddCounter(metric.ID, *metric.Delta)
		return err
	default:
		return ErrBadValue
	}
}
