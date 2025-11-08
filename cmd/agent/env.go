package main

import (
	"os"
	"strconv"
)

var (
	envRunAddr        string
	envReportInterval int64
	envPollInterval   int64
)

func parseEnv() {
	envRunAddr = os.Getenv("ADDRESS")

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
}
