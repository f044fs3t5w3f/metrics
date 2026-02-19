package configuration

import (
	"os"
	"strconv"
)

func EnvOrString(env string, fallback string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}
	return fallback
}

func EnvOrBool(env string, fallback bool) (bool, error) {
	if v := os.Getenv(env); v != "" {
		parsed, err := strconv.ParseBool(v)
		return parsed, err
	}
	return fallback, nil
}

func EnvOrInt64(env string, fallback int64) (int64, error) {
	if v := os.Getenv(env); v != "" {
		return strconv.ParseInt(v, 10, 64)
	}
	return fallback, nil
}
