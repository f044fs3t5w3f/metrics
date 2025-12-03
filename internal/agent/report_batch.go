package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/internal/sign"
	"go.uber.org/zap"
)

func ReportBatch(host string, batch MetricsBatch, key string) {
	url := fmt.Sprintf("http://%s/updates/", host)
	logger.Log.Info("to send metrics")

	logError := func(err error, attempt uint8) {
		logger.Log.Error("sendZippedJSON fail. Gonna retry", zap.Uint8("attempt", attempt), zap.Error(err))
	}

	// Маршалинг с очень маленькой вероятностью может дать ошибку в продакшене, поэтому нет ничего страшного,
	// что потенциально мы можем и эту ошибку ретраить, что казалось бы бесполезно и безнадёжно.
	err := retry(func() error {
		return sendZippedJSON(url, batch, key)
	}, []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}, logError)
	if err != nil {
		logger.Log.Error(err.Error())
	}
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

func sendZippedJSON(url string, data any, key string) error {
	body, err := getRequestBody(data)
	if err != nil {
		return fmt.Errorf("getRequestBody: %s", err)
	}
	var hash []byte
	if key != "" {
		signFunc := sign.GetSignFunc(key)
		data, _ := io.ReadAll(body)
		hash = signFunc(data)
		body = bytes.NewReader(data)
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("creating request error: %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	if len(hash) > 0 {
		base64String := base64.StdEncoding.EncodeToString(hash)
		req.Header.Set("HashSHA256", base64String)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("POST: %w", err)
	}
	response.Body.Close()
	return nil
}

// func reportMetric(host string, metric *models.Metrics) error {
// 	url := fmt.Sprintf("http://%s/update/", host)
// 	logger.Log.Info("to send metric", zap.String("type", metric.MType), zap.String("name", metric.ID))
// 	return sendZippedJSON(url, metric)
// }
