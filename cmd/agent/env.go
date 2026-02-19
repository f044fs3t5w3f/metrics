package main

import (
	"os"
	"strconv"
)

var (
	envRunAddr        string
	envReportInterval int64
	envPollInterval   int64
	envKey            string
	envRateLimit      int64
	envCryptoKeyPath  string
)

func parseEnv() {
	envRunAddr = os.Getenv("ADDRESS")
	envKey = os.Getenv("KEY")
	envCryptoKeyPath = os.Getenv("CRYPTO_KEY")

	var err error
	envReportIntervalStr := os.Getenv("REPORT_INTERVAL")
	if envReportIntervalStr != "" {
		envReportInterval, err = strconv.ParseInt(envReportIntervalStr, 10, 64)
		if err != nil {
			panic("REPORT_INTERVAL parse error")
		}
	}

	envPollIntervalStr := os.Getenv("POLL_INTERVAL")
	if envPollIntervalStr != "" {
		envPollInterval, err = strconv.ParseInt(envPollIntervalStr, 10, 64)
		if err != nil {
			panic("POLL_INTERVAL parse error")
		}
	}

	envRateLimitStr := os.Getenv("RATE_LIMIT")
	if envRateLimitStr != "" {
		envRateLimit, err = strconv.ParseInt(envRateLimitStr, 10, 64)
		if err != nil {
			panic("RATE_LIMIT parse error")
		}
	}

}
