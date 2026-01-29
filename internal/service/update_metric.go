package service

import (
	"context"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func (s *Service) UpdateMetric(ctx context.Context, metric models.Metrics) error {
	ev := audit.Event{
		Metrics:   []string{metric.ID},
		IP:        ctx.Value(CtxUserIP).(string),
		Timestamp: time.Now().Unix(),
	}
	s.audit.Notify(&ev)

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
