package null

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON is a json.RawMessage that will marshall as null when empty or nil.
type JSON json.RawMessage

var NullJSON = JSON(`null`)

// IsNull returns whether this JSON value is empty or contains null.
func (j JSON) IsNull() bool {
	return len(j) == 0 || bytes.Equal(j, NullJSON)
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
		*j = NullJSON
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
		*j = NullJSON
		return nil
	}

	if !json.Valid(raw) {
		return fmt.Errorf("scanned JSON isn't valid")
	}

	*j = raw
	return nil
}

func JSONValue(j JSON) (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return []byte(j), nil
}

func UnmarshalJSON(data []byte, j *JSON) error {
	return json.Unmarshal(data, (*json.RawMessage)(j))
}

func MarshalJSON(j JSON) ([]byte, error) {
	if len(j) == 0 {
		return json.Marshal(nil)
	}
	return []byte(j), nil
}
