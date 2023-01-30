package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"golang.org/x/exp/constraints"
)

// Int is an int that will write as null when it is zero both to databases and JSON
// null values when unmarshalled or scanned from a DB will result in a zero value.
type Int int

// NullInt is our constant for an Int value that will be written as null
const NullInt = Int(0)

// Scan implements the Scanner interface
func (i *Int) Scan(value any) error {
	return ScanInt(value, i)
}

// Value implements the Valuer interface
func (i Int) Value() (driver.Value, error) {
	return IntValue(i)
}

// UnmarshalJSON implements the Unmarshaller interface
func (i *Int) UnmarshalJSON(b []byte) error {
	return UnmarshalInt(b, i)
}

// MarshalJSON implements the Marshaller interface
func (i Int) MarshalJSON() ([]byte, error) {
	return MarshalInt(i)
}

// Int64 is an int64 that will write as null when it is zero both to databases and JSON
// null values when unmarshalled or scanned from a DB will result in a zero value.
type Int64 int64

// NullInt64 is our constant for an Int64 value that will be written as null
const NullInt64 = Int64(0)

// Scan implements the Scanner interface
func (i *Int64) Scan(value any) error {
	return ScanInt(value, i)
}

// Value implements the Valuer interface
func (i Int64) Value() (driver.Value, error) {
	return IntValue(i)
}

// UnmarshalJSON implements the Unmarshaller interface
func (i *Int64) UnmarshalJSON(b []byte) error {
	return UnmarshalInt(b, i)
}

// MarshalJSON implements the Marshaller interface
func (i Int64) MarshalJSON() ([]byte, error) {
	return MarshalInt(i)
}

// ScanInt scans a nullable INT into an int type, using zero for NULL.
func ScanInt[T constraints.Signed](value any, i *T) error {
	ni := sql.NullInt64{}

	if err := ni.Scan(value); err != nil {
		return err
	}

	if !ni.Valid {
		*i = T(0)
		return nil
	}

	*i = T(ni.Int64)
	return nil
}

// IntValue converts an int type value to NULL if it is zero.
func IntValue[T constraints.Signed](i T) (driver.Value, error) {
	if i == 0 {
		return nil, nil
	}
	return int64(i), nil
}

// UnmarshalInt unmarshals an int type from JSON, using zero for null.
func UnmarshalInt[T constraints.Signed](b []byte, i *T) error {
	var val *int64

	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}

	if val == nil {
		*i = 0
		return nil
	}

	*i = T(*val)
	return nil
}

// MarshalJSON marshals an int type to JSON, using null for zero.
func MarshalInt[T constraints.Signed](i T) ([]byte, error) {
	if i == 0 {
		return json.Marshal(nil)
	}
	return json.Marshal(int64(i))
}
