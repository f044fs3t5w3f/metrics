package config

import (
	"flag"

	"github.com/f044fs3t5w3f/metrics/pkg/configuration"
)

type Config struct {
	RunAddr        string
	ReportInterval int64
	PollInterval   int64
	Key            string
	RateLimit      int64
	CryptoKeyPath  string
}

func GetConfig() (*Config, error) {
	config := Config{}

	flagKey := flag.String("k", "", "key")
	flagEndpointAddr := flag.String("a", "localhost:8080", "endpoint address and port")
	flagCryptoPath := flag.String("crypto-key", "", "public key")
	flagReportInterval := flag.Int64("r", 10, "report interval")
	flagPollInterval := flag.Int64("p", 2, "poll interval")
	flagRateLimit := flag.Int64("l", 0, "rate limit")

	flag.Parse()

	config.Key = configuration.EnvOrString("KEY", *flagKey)
	config.RunAddr = configuration.EnvOrString("ADDRESS", *flagEndpointAddr)
	config.CryptoKeyPath = configuration.EnvOrString("CRYPTO_KEY", *flagCryptoPath)

	ReportInterval, err := configuration.EnvOrInt64("REPORT_INTERVAL", *flagReportInterval)
	if err != nil {
		return nil, err
	}
	config.ReportInterval = ReportInterval

	PollInterval, err := configuration.EnvOrInt64("POLL_INTERVAL", *flagPollInterval)
	if err != nil {
		return nil, err
	}
	config.PollInterval = PollInterval

	RateLimit, err := configuration.EnvOrInt64("REPORT_INTERVAL", *flagRateLimit)
	if err != nil {
		return nil, err
	}
	config.RateLimit = RateLimit

	return &config, nil
}
