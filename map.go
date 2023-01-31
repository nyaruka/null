package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Map is a generic map which is written to the database as JSON.
type Map map[string]any

// Scan implements the Scanner interface
func (m *Map) Scan(value any) error { return ScanMap(value, m) }

// Value implements the Valuer interface
func (m Map) Value() (driver.Value, error) { return MapValue(m) }

// UnmarshalJSON implements the Unmarshaller interface
func (m *Map) UnmarshalJSON(data []byte) error { return UnmarshalMap(data, m) }

// MarshalJSON implements the Marshaller interface
func (m Map) MarshalJSON() ([]byte, error) { return MarshalMap(m) }

// ScanMap scans a nullable text or JSON into a map, using an empty map for NULL.
func ScanMap(value any, m *Map) error {
	if value == nil {
		*m = make(Map)
		return nil
	}

	var raw []byte
	switch typed := value.(type) {
	case string:
		raw = []byte(typed)
	case []byte:
		raw = typed
	default:
		return fmt.Errorf("unable to scan %T as map", value)
	}

	// empty bytes is same as nil
	if len(raw) == 0 {
		*m = make(Map)
		return nil
	}

	if err := json.Unmarshal(raw, m); err != nil {
		return err
	}

	return nil
}

// MapValue converts a map to NULL if it is empty.
func MapValue(m Map) (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}

// MarshalMap marshals a map, returning null for an empty map.
func MarshalMap(m Map) ([]byte, error) {
	if len(m) == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(map[string]any(m))
}

func UnmarshalMap(data []byte, m *Map) error {
	err := json.Unmarshal(data, (*map[string]any)(m))
	if err != nil {
		return err
	}

	if *m == nil {
		*m = make(Map) // initialize empty map
	}
	return nil
}
