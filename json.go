package null

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSON is a json.RawMessage that will marshall as null when empty or nil.
// null and {} values when unmashalled or scanned from a DB will result in a nil value
type JSON json.RawMessage

// Scan implements the Scanner interface for decoding from a database
func (j *JSON) Scan(src interface{}) error {
	if src == nil {
		*j = nil
		return nil
	}

	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return fmt.Errorf("incompatible type for JSON type")
	}

	if !json.Valid(source) {
		return fmt.Errorf("invalid json: %s", source)
	}
	*j = source
	return nil
}

// Value implements the driver Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

// MarshalJSON encodes our JSON to JSON or null
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return json.Marshal(nil)
	}
	return []byte(j), nil
}

// UnmarshalJSON sets our JSON from the passed in JSON
func (j *JSON) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*j = nil
		return nil
	}

	var jj json.RawMessage
	err := json.Unmarshal(data, &jj)
	if err != nil {
		return err
	}

	*j = JSON(jj)
	return nil
}
