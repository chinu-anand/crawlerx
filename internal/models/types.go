package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type StringArray []string

// GORM requires Scanner and Valuer interfaces

func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("Scan failed: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}
