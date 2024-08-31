package main

import (
	"os"
	"strconv"
	"time"
)

func envString(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

func envInt64(key string, defaultValue int64) int64 {
	raw, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

func envDuration(key string, defaultValue time.Duration) time.Duration {
	raw, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return defaultValue
	}

	return value
}
