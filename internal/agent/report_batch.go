package agent

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/pkg/net"
	"github.com/f044fs3t5w3f/metrics/pkg/retry"
	"go.uber.org/zap"
)

func ReportBatch(host string, batch MetricsBatch, key string, publicKey *rsa.PublicKey) {
	url := fmt.Sprintf("http://%s/updates/", host)
	logger.Log.Info("to send metrics")

	logError := func(err error, attempt uint8) {
		logger.Log.Error("sendZippedJSON fail. Gonna retry", zap.Uint8("attempt", attempt), zap.Error(err))
	}

	// Маршалинг с очень маленькой вероятностью может дать ошибку в продакшене, поэтому нет ничего страшного,
	// что потенциально мы можем и эту ошибку ретраить, что казалось бы бесполезно и безнадёжно.

	err := retry.Retry(func() error {
		return net.SendZippedSignedJSON(url, batch, key, publicKey)
	}, []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}, logError)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
