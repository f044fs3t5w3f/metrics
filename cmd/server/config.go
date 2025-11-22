package main

import (
	"flag"
	"os"
	"strconv"
)

type config struct {
	runAddr         string
	fileStoragePath string
	restore         bool
	storeInterval   int64
	databaseParams  string
}

func getString(envName string, flagName string, value string, usage string) string {
	envValue := os.Getenv(envName)
	if envValue != "" {
		return envValue
	}
	return *flag.String(flagName, value, usage)
}

func getInt64(envName string, flagName string, value int64, usage string) (int64, error) {
	envValue := os.Getenv(envName)
	if envValue != "" {
		return strconv.ParseInt(envValue, 10, 64)
	}
	return *flag.Int64(flagName, value, usage), nil
}

func getBool(envName string, flagName string, value bool, usage string) (bool, error) {
	envValue := os.Getenv(envName)
	if envValue != "" {
		return strconv.ParseBool(envValue)
	}
	return *flag.Bool(flagName, value, usage), nil
}

func getConfig() (*config, error) {
	flag.Parse()
	config := config{}
	config.runAddr = getString("ADDRESS", "a", "localhost:8080", "address and port to run server")
	config.fileStoragePath = getString("FILE_STORAGE_PATH", "f", "store.json", "storage file name. Default store.json")
	config.databaseParams = getString("DATABASE_DSN", "d", "", "Database params")
	storeInterval, err := getInt64("STORE_INTERVAL", "i", 300, "interval (seconds) for saving metrics to disk. 0 enables synchronous writes. Default: 300")
	if err != nil {
		return nil, err
	}
	config.storeInterval = storeInterval

	restore, err := getBool("RESTORE", "r", false, "restore metrics from file")
	if err != nil {
		return nil, err
	}
	config.restore = restore

	return &config, nil
}

// getEnvString(envName string, fallBack string) err
