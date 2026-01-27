package service

import (
	"context"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func (s *Service) UpdateMetrics(ctx context.Context, metrics []*models.Metrics) error {
	metricNames := make([]string, len(metrics))
	for i, metric := range metrics {
		metricNames[i] = metric.ID
	}
	ev := audit.Event{
		Metrics:   metricNames,
		IP:        ctx.Value(CtxUserIP).(string),
		Timestamp: time.Now().Unix(),
	}
	s.audit.Notify(&ev)
	err := s.storage.MultiUpdate(metrics)
	return err
}
