package main

import (
	"flag"
)

var (
	flagEndpointAddr   string
	flagReportInterval int64
	flagPollInterval   int64
)

func parseFlags() {
	flag.StringVar(&flagEndpointAddr, "a", "localhost:8080", "endpoint address and port")
	flag.Int64Var(&flagReportInterval, "r", 10, "report interval")
	flag.Int64Var(&flagPollInterval, "p", 2, "poll interval")
	flag.Parse()
}
