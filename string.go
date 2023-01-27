package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

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
