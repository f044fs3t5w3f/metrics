package main

import (
	"os"
)

var envRunAddr string

func parseEnv() {
	envRunAddr = os.Getenv("ADDRESS")
}
