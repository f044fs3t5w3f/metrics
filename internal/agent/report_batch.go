package agent

import (
	"fmt"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

func ReportBatch(host string, batch MetricsBatch) {
	for _, metric := range batch {
		reportMetric(host, metric)
	}
}

func reportMetric(host string, metric *models.Metrics) error {
	url, err := getURL(host, metric)
	if err != nil {
		return fmt.Errorf("getURL: %s", err)
	}
	response, err := http.Post(url, "", nil)
	if err != nil {
		return fmt.Errorf("POST: %s", err)
	}
	response.Body.Close()
	return nil
}
