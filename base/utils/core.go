package utils

import (
	"os"
	"strconv"
)

// DefValueType value types for default value of environment variables string | int | bool
type DefValueType interface {
	string | int | bool
}

// GetEnv Load environment variable or return default value of DefValueType
func GetEnv[T DefValueType](key string, def T) T {
	var ret T
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}

	// switch on the pointer types of T
	switch p := any(&ret).(type) {
	case *string:
		*p = v
	case *int:
		iv, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		*p = iv
	case *bool:
		bv, err := strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
		*p = bv
	}
	return ret
}

func CopyMap[K comparable, V any](src map[K]V, dst map[K]V) map[K]V {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
