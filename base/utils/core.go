package utils

import (
	"os"
)

// Getenv returns enviroment variable value.
// If variable does not exist, returns default value.
func Getenv(key, def string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return def
}
