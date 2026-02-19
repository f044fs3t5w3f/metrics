package main

import (
	"flag"

	"github.com/f044fs3t5w3f/metrics/pkg/configuration"
)

type config struct {
	runAddr         string
	fileStoragePath string
	restore         bool
	storeInterval   int64
	databaseParams  string
	key             string
	auditURL        string
	auditFile       string
	cryptoFile      string
}

func getConfig() (*config, error) {
	config := &config{}

	addrFlag := flag.String("a", "localhost:8080", "server address")
	fileFlag := flag.String("f", "store.json", "storage file")
	keyFlag := flag.String("k", "", "key")
	dbFlag := flag.String("d", "", "database dsn")
	intervalFlag := flag.Int64("i", 300, "store interval")
	restoreFlag := flag.Bool("r", false, "restore on startup")
	auditFile := flag.String("audit-file", "", "audit file")
	auditURL := flag.String("audit-url", "", "audit url")
	cryptoFile := flag.String("crypto-key", "", "private key")

	flag.Parse()

	config.runAddr = configuration.EnvOrString("ADDRESS", *addrFlag)
	config.fileStoragePath = configuration.EnvOrString("FILE_STORAGE_PATH", *fileFlag)
	config.databaseParams = configuration.EnvOrString("DATABASE_DSN", *dbFlag)
	config.key = configuration.EnvOrString("KEY", *keyFlag)
	config.auditFile = configuration.EnvOrString("AUDIT_FILE", *auditFile)
	config.auditURL = configuration.EnvOrString("AUDIT_URL", *auditURL)
	config.cryptoFile = configuration.EnvOrString("CRYPTO_KEY", *cryptoFile)

	interval, err := configuration.EnvOrInt64("STORE_INTERVAL", *intervalFlag)
	if err != nil {
		return nil, err
	}
	config.storeInterval = interval

	restore, err := configuration.EnvOrBool("RESTORE", *restoreFlag)
	if err != nil {
		return nil, err
	}
	config.restore = restore

	return config, nil
}
