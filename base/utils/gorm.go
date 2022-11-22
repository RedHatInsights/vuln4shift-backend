package utils

import (
	"database/sql/driver"
)

// ByteArrayBool sets to true if the array is not null or empty during the GORM scan, else false.
type ByteArrayBool bool

func (s *ByteArrayBool) Scan(value interface{}) error {
	if value == nil {
		*s = false
	}

	res := value.([]byte)
	if len(res) > 0 {
		*s = true
	} else {
		*s = false
	}

	return nil
}

func (s ByteArrayBool) Value() (driver.Value, error) {
	if s {
		return "true", nil
	}
	return "false", nil
}
