package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// String is string that will write as null when it is empty, both to databases and JSON
// null values when unmarshalled or scanned from a DB will result in an empty string value.
type String string

// NullString is our constant for an String value that will be written as null
const NullString = String("")

// Scan implements the Scanner interface
func (s *String) Scan(value any) error {
	return ScanString(value, s)
}

// Value implements the Valuer interface
func (s String) Value() (driver.Value, error) {
	return StringValue(s)
}

// UnmarshalJSON implements the Unmarshaller interface
func (s *String) UnmarshalJSON(b []byte) error {
	return UnmarshalString(b, s)
}

// MarshalJSON implements the Marshaller interface
func (s String) MarshalJSON() ([]byte, error) {
	return MarshalString(s)
}

// ScanString scans a nullable CHAR/TEXT into a string type, using empty string for NULL.
func ScanString[T ~string](value any, s *T) error {
	ns := sql.NullString{}

	if err := ns.Scan(value); err != nil {
		return err
	}

	if !ns.Valid {
		*s = ""
		return nil
	}

	*s = T(ns.String)
	return nil
}

// StringValue converts a string type value to NULL if it is empty
func StringValue[T ~string](s T) (driver.Value, error) {
	if s == "" {
		return nil, nil
	}
	return string(s), nil
}

// UnmarshalString unmarshals a string type from JSON, using empty string for null.
func UnmarshalString[T ~string](b []byte, s *T) error {
	var val *string

	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	if val == nil {
		*s = ""
		return nil
	}

	*s = T(*val)
	return nil
}

// MarshalJSON marshals a string type to JSON, using null for empty strings.
func MarshalString[T ~string](s T) ([]byte, error) {
	if s == "" {
		return json.Marshal(nil)
	}
	return json.Marshal(string(s))
}
