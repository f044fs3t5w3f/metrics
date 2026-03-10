package main

type Config struct {
	Address        string `env:"ADDRESS" flag:"a" jsonConfig:"address" default:"localhost:8080" description:"endpoint address and port"`
	ReportInterval int64  `env:"REPORT_INTERVAL" flag:"r" jsonConfig:"report_interval" default:"10" description:"report interval in seconds"`
	PollInterval   int64  `env:"POLL_INTERVAL" flag:"p" jsonConfig:"poll_interval" default:"1" description:"poll interval in seconds"`
	Key            string `env:"KEY" flag:"k" jsonConfig:"key" default:"" description:"key"`
	CryptoKeyPath  string `env:"CRYPTO_KEY" flag:"crypto-key" jsonConfig:"crypto_key" default:"" description:"public key path"`
	RateLimit      int64  `env:"RATE_LIMIT" flag:"l" jsonConfig:"rate_limit" default:"0" description:"rate limit"`
	RPCServer      string `env:"RPC_SERVER" flag:"rpc" jsonConfig:"rpc_server" default:"" description:"rpc server address"`
}
