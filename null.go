package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// Int is an int that will write as null when it is zero both to databases and json
// null values when unmashalled or scanned from a DB will result in a 0 value
type Int int64

// NullInt is our constant for an Int value that will be written as null
const NullInt = Int(0)

// UnmarshalInt is a utility method that can be used to unmarshal a json value to an Int and error
// In the case of an error, null or "" values, NullInt is returned
func UnmarshalInt(b []byte) (Int, error) {
	var val *int64
	err := json.Unmarshal(b, &val)
	if err != nil || val == nil {
		return NullInt, err
	}
	return Int(*val), nil
}

// ScanInt is a utility method that can be used to scan a db value and return an Int and error
// In the case of an error, null or "" values, NullInt is returned
func ScanInt(value interface{}) (Int, error) {
	val := &sql.NullInt64{}
	err := val.Scan(value)
	if err != nil {
		return NullInt, err
	}

	if !val.Valid {
		return NullInt, nil
	}

	return Int(val.Int64), nil
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
	val, err := UnmarshalInt(b)
	*i = val
	return err
}

// Scan implements the Scanner interface for Int
func (i *Int) Scan(value interface{}) error {
	val, err := ScanInt(value)
	*i = val
	return err
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
func UnmarshalString(b []byte) (String, error) {
	var val *string
	err := json.Unmarshal(b, &val)
	if err != nil {
		return NullString, err
	}
	if val == nil {
		return NullString, nil
	}

	return String(*val), nil
}

// ScanString is a utility method that can be used to scan a db value and return a String and error
// In the case of an error, null or "" values, NullString is returned
func ScanString(value interface{}) (String, error) {
	val := &sql.NullString{}
	err := val.Scan(value)
	if err != nil {
		return NullString, err
	}

	if !val.Valid {
		return NullString, nil
	}

	return String(val.String), nil
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
	val, err := UnmarshalString(b)
	*s = val
	return err
}

// Scan implements the Scanner interface for String
func (s *String) Scan(value interface{}) error {
	val, err := ScanString(value)
	*s = val
	return err
}

// Value implements the driver Valuer interface for String
func (s String) Value() (driver.Value, error) {
	if s == NullString {
		return nil, nil
	}
	return string(s), nil
}
