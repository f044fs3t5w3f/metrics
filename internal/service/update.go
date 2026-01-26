package service

import (
	"errors"
	"strconv"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func (s *Service) Update(type_, metricName, merticValueStr string) error {
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
