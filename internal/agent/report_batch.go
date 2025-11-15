package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/models"
	"go.uber.org/zap"
)

func ReportBatch(host string, batch MetricsBatch) {
	url := fmt.Sprintf("http://%s/updates/", host)
	logger.Log.Info("to send metrics")
	err := sendZippedJSON(url, batch)
	if err != nil {
		logger.Log.Error(err.Error())
	}
	// for _, metric := range batch {
	// 	err := reportMetric(host, metric)
	// 	if err != nil {
	// 		logger.Log.Error(err.Error())
	// 	}
	// }
}

func getRequestBody(data any) (io.Reader, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshalling metric error: %s", err)
	}
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write(jsonData)
	gz.Close()
	return &buf, nil
}

func sendZippedJSON(url string, data any) error {
	body, err := getRequestBody(data)
	if err != nil {
		return fmt.Errorf("getRequestBody: %s", err)
	}
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("creating request error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST: %s", err)
	}
	response.Body.Close()
	return nil
}

func reportMetric(host string, metric *models.Metrics) error {
	url := fmt.Sprintf("http://%s/update/", host)
	logger.Log.Info("to send metric", zap.String("type", metric.MType), zap.String("name", metric.ID))
	return sendZippedJSON(url, metric)
}
