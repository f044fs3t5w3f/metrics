package service

import "github.com/f044fs3t5w3f/metrics/internal/models"

func (s *Service) UpdateMetrics(metrics []*models.Metrics) error {
	err := s.storage.MultiUpdate(metrics)
	return err
}
