package agent

import (
	"bytes"
	"encoding/json"
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
	url := fmt.Sprintf("http://%s/update/", host)
	logger.Log.Info("to send metric", zap.String("type", metric.MType), zap.String("name", metric.ID))
	jsonData, _ := json.Marshal(metric)
	body := bytes.NewBuffer(jsonData)
	response, err := http.Post(url, "", body)
	if err != nil {
		return fmt.Errorf("POST: %s", err)
	}
	response.Body.Close()
	return nil
}
