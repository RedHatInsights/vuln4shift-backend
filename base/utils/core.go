package utils

import (
	"os"
	"strconv"
)

// GetEnv Load string environment variable or return default value
func GetEnv(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

// GetIntEnv Load integer environment variable or return default value
func GetIntEnv(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	parsedInt, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}

	return parsedInt
}
