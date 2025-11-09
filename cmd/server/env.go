package main

import (
	"os"
)

var (
	envRunAddr, envFileStoragePath, envStoreInterval, envRestore, envDatabaseParams string
)

func parseEnv() {
	envRunAddr = os.Getenv("ADDRESS")
	envFileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	envStoreInterval = os.Getenv("STORE_INTERVAL")
	envRestore = os.Getenv("RESTORE")
	envDatabaseParams = os.Getenv("DATABASE_DSN")
}
