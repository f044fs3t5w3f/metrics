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
	key             string
}

func envOrString(env string, fallback string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return fallback
}

func envOrBool(env string, fallback bool) (bool, error) {
	if v := os.Getenv(env); v != "" {
		parsed, err := strconv.ParseBool(v)
		return parsed, err
	}
	return fallback, nil
}

func envOrInt64(env string, fallback int64) (int64, error) {
	if v := os.Getenv(env); v != "" {
		return strconv.ParseInt(v, 10, 64)
	}
	return fallback, nil
}

func getConfig() (*config, error) {
	config := &config{}

	addrFlag := flag.String("a", "localhost:8080", "server address")
	fileFlag := flag.String("f", "store.json", "storage file")
	keyFlag := flag.String("k", "", "key")
	dbFlag := flag.String("d", "", "database dsn")
	intervalFlag := flag.Int64("i", 300, "store interval")
	restoreFlag := flag.Bool("r", false, "restore on startup")

	flag.Parse()

	config.runAddr = envOrString("ADDRESS", *addrFlag)
	config.fileStoragePath = envOrString("FILE_STORAGE_PATH", *fileFlag)
	config.databaseParams = envOrString("DATABASE_DSN", *dbFlag)
	config.key = envOrString("KEY", *keyFlag)

	interval, err := envOrInt64("STORE_INTERVAL", *intervalFlag)
	if err != nil {
		return nil, err
	}
	config.storeInterval = interval

	restore, err := envOrBool("RESTORE", *restoreFlag)
	if err != nil {
		return nil, err
	}
	config.restore = restore

	return config, nil
}
