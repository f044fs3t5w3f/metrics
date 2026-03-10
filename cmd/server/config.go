package main

type Config struct {
	RunAddr       string `env:"ADDRESS" flag:"a" jsonConfig:"address" default:"localhost:8080" description:"Server address"`
	FileStorage   string `env:"FILE_STORAGE_PATH" flag:"f" jsonConfig:"store_file" default:"store.json" description:"Storage file path"`
	Key           string `env:"KEY" flag:"k" jsonConfig:"crypto_key" default:"" description:"Key for encryption purposes"`
	AuditFile     string `env:"AUDIT_FILE" flag:"audit-file" jsonConfig:"audit_file" default:"" description:"Audit log file path"`
	AuditURL      string `env:"AUDIT_URL" flag:"audit-url" jsonConfig:"audit_url" default:"" description:"Audit service URL"`
	CryptoFile    string `env:"CRYPTO_KEY" flag:"crypto-key" jsonConfig:"crypto_key_path" default:"" description:"Private crypto key file path"`
	RestoreOnBoot bool   `env:"RESTORE" flag:"r" jsonConfig:"restore" default:"false" description:"Restore data from storage on bootup"`
	StoreInterval int64  `env:"STORE_INTERVAL" flag:"i" jsonConfig:"store_interval" default:"300" description:"Interval in seconds to store metrics"`
	DatabaseDSN   string `env:"DATABASE_DSN" flag:"d" jsonConfig:"database_dsn" default:"" description:"Database connection string"`
	TrustedSubnet string `env:"TRUSTED_SUBNET" flag:"t" jsonConfig:"trusted_subnet" default:""`
	RPCServer     string `env:"RPC_SERVER" flag:"rpc" jsonConfig:"rpc_server" default:"" description:"rpc server address"`
}

// func getConfig() (*config, error) {
// 	config := &config{}

// 	addrFlag := flag.String("a", "localhost:8080", "server address")
// 	fileFlag := flag.String("f", "store.json", "storage file")
// 	keyFlag := flag.String("k", "", "key")
// 	dbFlag := flag.String("d", "", "database dsn")
// 	intervalFlag := flag.Int64("i", 300, "store interval")
// 	restoreFlag := flag.Bool("r", false, "restore on startup")
// 	auditFile := flag.String("audit-file", "", "audit file")
// 	auditURL := flag.String("audit-url", "", "audit url")
// 	cryptoFile := flag.String("crypto-key", "", "private key")

// 	flag.Parse()

// 	config.runAddr = configuration.EnvOrString("ADDRESS", *addrFlag)
// 	config.fileStoragePath = configuration.EnvOrString("FILE_STORAGE_PATH", *fileFlag)
// 	config.databaseParams = configuration.EnvOrString("DATABASE_DSN", *dbFlag)
// 	config.key = configuration.EnvOrString("KEY", *keyFlag)
// 	config.auditFile = configuration.EnvOrString("AUDIT_FILE", *auditFile)
// 	config.auditURL = configuration.EnvOrString("AUDIT_URL", *auditURL)
// 	config.cryptoFile = configuration.EnvOrString("CRYPTO_KEY", *cryptoFile)

// 	interval, err := configuration.EnvOrInt64("STORE_INTERVAL", *intervalFlag)
// 	if err != nil {
// 		return nil, err
// 	}
// 	config.storeInterval = interval

// 	restore, err := configuration.EnvOrBool("RESTORE", *restoreFlag)
// 	if err != nil {
// 		return nil, err
// 	}
// 	config.restore = restore

// 	return config, nil
// }
