package agent

import (
	"fmt"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
	"go.uber.org/zap"
)

func ReportBatch(host string, batch MetricsBatch) {
	for _, metric := range batch {
		err := reportMetric(host, metric)
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}
}

func reportMetric(host string, metric *models.Metrics) error {
	url, err := getURL(host, metric)
	if err != nil {
		return fmt.Errorf("getURL: %s", err)
	}
	logger.Log.Info("to send metric", zap.String("url", url))
	response, err := http.Post(url, "", nil)
	if err != nil {
		return fmt.Errorf("POST: %s", err)
	}
	response.Body.Close()
	return nil
}
