package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/audit"
	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func (s *Service) Update(ctx context.Context, type_, metricName, merticValueStr string) error {
	ev := audit.Event{
		Metrics:   []string{metricName},
		IP:        ctx.Value(CtxUserIP).(string),
		Timestamp: time.Now().Unix(),
	}
	s.audit.Notify(&ev)
	switch type_ {
	case models.Gauge:
		merticParsed, err := strconv.ParseFloat(merticValueStr, 64)
		if err != nil {
			return err
		}
		s.storage.SetGauge(metricName, merticParsed)
	case models.Counter:
		merticParsed, err := strconv.ParseInt(merticValueStr, 0, 64)
		if err != nil {
			return err
		}
		s.storage.AddCounter(metricName, merticParsed)
	default:
		return errors.New("unknown type")
	}
	return nil
}
