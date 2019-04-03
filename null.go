package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Int is an int that will write as null when it is zero both to databases and json
// null values when unmashalled or scanned from a DB will result in a 0 value
type Int int64

// NullInt is our constant for an Int value that will be written as null
const NullInt = Int(0)

// UnmarshalInt is a utility method that can be used to unmarshal a json value to an Int and error
// In the case of an error, null or "" values, NullInt is returned
func UnmarshalInt(b []byte, v *Int) error {
	val := int64(0)
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	*v = Int(val)
	return nil
}

// ScanInt is a utility method that can be used to scan a db value and return an Int and error
// In the case of an error, null or "" values, NullInt is returned
func ScanInt(value interface{}, v *Int) error {
	val := &sql.NullInt64{}
	err := val.Scan(value)
	if err != nil {
		return err
	}

	if !val.Valid {
		*v = NullInt
		return nil
	}

	*v = Int(val.Int64)
	return nil
}

// MarshalJSON marshals our int to JSON. 0 values will be marshalled as null
func (i Int) MarshalJSON() ([]byte, error) {
	if i == NullInt {
		return json.Marshal(nil)
	}
	return json.Marshal(int64(i))
}

// UnmarshalJSON unmarshals our JSON to int. null values will be marshalled to 0
func (i *Int) UnmarshalJSON(b []byte) error {
	return UnmarshalInt(b, i)
}

// Scan implements the Scanner interface for Int
func (i *Int) Scan(value interface{}) error {
	return ScanInt(value, i)
}

// Value implements the driver Valuer interface for Int
func (i Int) Value() (driver.Value, error) {
	if i == NullInt {
		return nil, nil
	}
	return int64(i), nil
}

// String is string that will write as null when it is "" both to databases and json
// null values when unmashalled or scanned from a DB will result in a "" value
type String string

// NullString is our constant for an String value that will be written as null
const NullString = String("")

// UnmarshalString is a utility method that can be used to unmarshal a json value to a String and error
// In the case of an error, null or "" values, NullString is returned
func UnmarshalString(b []byte, v *String) error {
	var val *string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return err
	}
	if val == nil {
		*v = NullString
		return nil
	}

	*v = String(*val)
	return nil
}

// ScanString is a utility method that can be used to scan a db value and return a String and error
// In the case of an error, null or "" values, NullString is returned
func ScanString(value interface{}, v *String) error {
	val := &sql.NullString{}
	err := val.Scan(value)
	if err != nil {
		return err
	}

	if !val.Valid {
		*v = NullString
		return nil
	}

	*v = String(val.String)
	return nil
}

// MarshalJSON marshals our string to JSON. "" values will be marshalled as null
func (s String) MarshalJSON() ([]byte, error) {
	if s == NullString {
		return json.Marshal(nil)
	}
	return json.Marshal(string(s))
}

// UnmarshalJSON unmarshals our json to a string. null values will be marshalled to ""
func (s *String) UnmarshalJSON(b []byte) error {
	return UnmarshalString(b, s)
}

// Scan implements the Scanner interface for String
func (s *String) Scan(value interface{}) error {
	return ScanString(value, s)
}

// Value implements the driver Valuer interface for String
func (s String) Value() (driver.Value, error) {
	if s == NullString {
		return nil, nil
	}
	return string(s), nil
}

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
