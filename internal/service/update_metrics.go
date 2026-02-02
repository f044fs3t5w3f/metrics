package service

import (
	"context"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func (s *Service) UpdateMetrics(ctx context.Context, metrics []*models.Metrics) error {
	metricNames := make([]string, len(metrics))
	for i, metric := range metrics {
		metricNames[i] = metric.ID
	}
	ip, ok := ctx.Value(CtxUserIP).(string)
	if !ok {
		logger.Log.Warn("incorrect CtxUserIP value cast to string")
	}
	ev := audit.Event{
		Metrics:   metricNames,
		IP:        ip,
		Timestamp: time.Now().Unix(),
	}
	s.audit.Notify(&ev)
	err := s.storage.MultiUpdate(metrics)
	return err
}
