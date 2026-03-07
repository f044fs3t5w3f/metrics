package agent

import (
	"crypto/rsa"
	"fmt"
	"sync"
	"time"

	"github.com/f044fs3t5w3f/metrics/internal/logger"
	"github.com/f044fs3t5w3f/metrics/pkg/net"
	"github.com/f044fs3t5w3f/metrics/pkg/retry"
	"go.uber.org/zap"
)

func ReportBatch(host string, batch MetricsBatch, key string, publicKey *rsa.PublicKey, wg *sync.WaitGroup) {
	url := fmt.Sprintf("http://%s/updates/", host)
	logger.Log.Info("to send metrics")

	logError := func(err error, attempt uint8) {
		logger.Log.Error("sendZippedJSON fail. Gonna retry", zap.Uint8("attempt", attempt), zap.Error(err))
	}

	// Маршалинг с очень маленькой вероятностью может дать ошибку в продакшене, поэтому нет ничего страшного,
	// что потенциально мы можем и эту ошибку ретраить, что казалось бы бесполезно и безнадёжно.

	err := retry.Retry(func() error {
		// Обычно wg.Add не должен вызываться в горутинах, но в этому случае это оправдано, так как wg используется
		// только для того, чтобы не прервать процесс приложения во время непостредствеено запроса к серверу
		// Вообще тут можно было обойтись ещё одним мьютексом, так как у нас нет такого, что параллельно выполняются
		// несколько запросов, которые мы должны дождаться. Но будем считать, что они могут появиться.
		wg.Add(1)
		err := net.SendZippedSignedJSON(url, batch, key, publicKey)
		wg.Done()
		return err
	}, []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}, logError)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
