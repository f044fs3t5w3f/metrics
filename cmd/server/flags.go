package main

import (
	"flag"
)

var (
	flagRunAddr         string
	flagFileStoragePath string
	flagRestore         bool
	flagStoreInterval   int64
	flagDatabaseParams  string
)

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&flagFileStoragePath, "f", "store.json", "storage file name. Default store.json")
	flag.Int64Var(&flagStoreInterval, "i", 300, "interval (seconds) for saving metrics to disk. 0 enables synchronous writes. Default: 300")
	flag.BoolVar(&flagRestore, "r", false, "restore metrics from file")
	flag.StringVar(&flagDatabaseParams, "d", "", "Database params")
	flag.Parse()
}
