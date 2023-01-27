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
