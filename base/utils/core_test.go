package utils

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvInvalidBool(t *testing.T) {
	assert.Panics(t, func() { GetEnv("LOGGING_LEVEL", false) })
}

func TestGetEnvInvalidInt(t *testing.T) {
	assert.Panics(t, func() { GetEnv("LOGGING_LEVEL", 0) })
}

func TestGetEnvBool(t *testing.T) {
	assert.Nil(t, os.Setenv("TEST_BOOLEAN", "TRUE"))
	assert.True(t, GetEnv("TEST_BOOLEAN", false))
}

func TestCopyMap(t *testing.T) {
	src := map[string]bool{
		"t": true,
		"f": false,
	}
	dst := make(map[string]bool)
	CopyMap(src, dst)
	assert.True(t, reflect.DeepEqual(src, dst))
}
