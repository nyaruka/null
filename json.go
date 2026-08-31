package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON is a json.RawMessage that will marshall as null when empty or nil.
type JSON json.RawMessage

// the null literal, kept as an untyped constant so that it can't be aliased or mutated
const nullJSON = "null"

// NullJSON is our constant for a JSON value that will be written as null. It must not be mutated or reassigned -
// use JSON(nil) or JSON(`null`) to construct a null value.
var NullJSON = JSON(nullJSON)

// IsNull returns whether this JSON value is empty or contains null.
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == nullJSON
}

// Scan implements the Scanner interface
func (j *JSON) Scan(value any) error { return ScanJSON(value, j) }

// Value implements the Valuer interface
func (j JSON) Value() (driver.Value, error) { return JSONValue(j) }

// UnmarshalJSON implements the Unmarshaller interface
func (j *JSON) UnmarshalJSON(data []byte) error { return UnmarshalJSON(data, j) }

// MarshalJSON implements the Marshaller interface
func (j JSON) MarshalJSON() ([]byte, error) { return MarshalJSON(j) }

func ScanJSON(value any, j *JSON) error {
	if value == nil {
		*j = JSON(nullJSON)
		return nil
	}

	var raw []byte
	switch typed := value.(type) {
	case string:
		raw = []byte(typed)
	case []byte:
		raw = typed
	default:
		return fmt.Errorf("unable to scan %T as JSON", value)
	}

	// empty bytes is same as nil
	if len(raw) == 0 {
		*j = JSON(nullJSON)
		return nil
	}

	if !json.Valid(raw) {
		return fmt.Errorf("scanned JSON isn't valid")
	}

	// we need to make our own copy of this data as the driver is allowed to reuse it for subsequent scans - usually
	// database/sql takes care of this but we're implementing our own Scanner https://github.com/golang/go/issues/24492
	*j = JSON(bytes.Clone(raw))
	return nil
}

func JSONValue(j JSON) (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return []byte(j), nil
}

func UnmarshalJSON(data []byte, j *JSON) error {
	// unmarshal into a new value rather than *j, because json.RawMessage appends into whatever backing array it
	// already has, which would write through to any value sharing it
	var raw json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*j = JSON(raw)
	return nil
}

func MarshalJSON(j JSON) ([]byte, error) {
	if len(j) == 0 {
		return json.Marshal(nil)
	}
	// return a copy so that callers can't mutate the value via the returned slice
	return bytes.Clone(j), nil
}
