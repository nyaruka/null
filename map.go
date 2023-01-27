package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Map is a one level deep dictionary that is represented as JSON text in the database.
// Empty maps will be written as null to the database and to JSON.
type Map struct {
	m map[string]interface{}
}

// NewMap creates a new Map
func NewMap(m map[string]interface{}) Map {
	return Map{m: m}
}

// Map returns our underlying map
func (m *Map) Map() map[string]interface{} {
	if m.m == nil {
		m.m = make(map[string]interface{})
	}
	return m.m
}

// GetString returns the string value with the passed in key, or def if not found or of wrong type
func (m *Map) GetString(key string, def string) string {
	if m.m == nil {
		return def
	}
	val := m.m[key]
	if val == nil {
		return def
	}
	str, isStr := val.(string)
	if !isStr {
		return def
	}
	return str
}

// Get returns the  value with the passed in key, or def if not found
func (m *Map) Get(key string, def interface{}) interface{} {
	if m.m == nil {
		return def
	}
	val := m.m[key]
	if val == nil {
		return def
	}
	return val
}

// Scan implements the Scanner interface for decoding from a database
func (m *Map) Scan(src interface{}) error {
	m.m = make(map[string]interface{})
	if src == nil {
		return nil
	}

	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return fmt.Errorf("incompatible type for map")
	}

	// 0 length string is same as nil
	if len(source) == 0 {
		return nil
	}

	err := json.Unmarshal(source, &m.m)
	if err != nil {
		return err
	}
	return nil
}

// Value implements the driver Valuer interface
func (m Map) Value() (driver.Value, error) {
	if m.m == nil || len(m.m) == 0 {
		return nil, nil
	}
	return json.Marshal(m.m)
}

// MarshalJSON encodes our map to JSON
func (m Map) MarshalJSON() ([]byte, error) {
	if m.m == nil || len(m.m) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(m.m)
}

// UnmarshalJSON sets our map from the passed in JSON
func (m *Map) UnmarshalJSON(data []byte) error {
	m.m = make(map[string]interface{})
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &m.m)
}
