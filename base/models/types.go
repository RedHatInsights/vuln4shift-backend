package models

import (
	"database/sql/driver"
	"errors"
)

// Severity represents severity enum from database
type Severity string

// Severity values in database
const (
	NotSet    Severity = "NotSet"
	None      Severity = "None"
	Low       Severity = "Low"
	Medium    Severity = "Medium"
	Moderate  Severity = "Moderate"
	Important Severity = "Important"
	High      Severity = "High"
	Critical  Severity = "Critical"
)

// Scan scanner interface implementation for Severity
func (s *Severity) Scan(value interface{}) error {
	if value == nil {
		return errors.New("invalid scan value for severity type")
	}
	*s = Severity(value.(string))
	return nil
}

// Value valuer interface implementation for Severity
func (s Severity) Value() (driver.Value, error) {
	return string(s), nil
}
