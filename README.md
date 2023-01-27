# null [![Build Status](https://github.com/nyaruka/null/workflows/CI/badge.svg)](https://github.com/nyaruka/null/actions?query=workflow%3ACI) [![codecov](https://codecov.io/gh/nyaruka/null/branch/main/graph/badge.svg)](https://codecov.io/gh/nyaruka/null) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/null)](https://goreportcard.com/report/github.com/nyaruka/null)

This module provides (yet another) alternative to dealing with null values in databases or JSON. Other approaches like 
the [Null types](https://golang.org/pkg/database/sql/#NullInt64) in the standard library use structs to ensure you can 
differentiate between zero values an null values. If that isn't a meaningful distinction in your app, then this module
might be a simpler approach for you because it uses primitive values and treats zero values as null values.

If you don't need to define your own types, you can use one of the following predefined types:

```go
null.Int    // 0 saves as NULL, NULL scans as zero
null.String // "" saves as NULL, NULL scans as ""
null.Map    // empty map saves as NULL, NULL scans as empty map
```

If you want to define a custom integer type, you need to define the following methods:

```go
import "github.com/nyaruka/null"

type CustomID null.Int

const NullCustomID = CustomID(0)

// MarshalJSON marshals into JSON. 0 values will become null
func (i CustomID) MarshalJSON() ([]byte, error) {
	return null.Int(i).MarshalJSON()
}

// UnmarshalJSON unmarshals from JSON. null values become 0
func (i *CustomID) UnmarshalJSON(b []byte) error {
	return null.UnmarshalInt(b, (*null.Int)(i))
}

// Value returns the db value, null is returned for 0
func (i CustomID) Value() (driver.Value, error) {
	return null.Int(i).Value()
}

// Scan scans from the db value. null values become 0
func (i *CustomID) Scan(value interface{}) error {
	return null.ScanInt(value, (*null.Int)(i))
}
```

And likewise for a custom string type:

```go
import "github.com/nyaruka/null"

type CustomString null.String

type NullCustomString = CustomString("")

// MarshalJSON marshals into JSON. "" values will become null
func (s CustomString) MarshalJSON() ([]byte, error) {
	return null.String(s).MarshalJSON()
}

// UnmarshalJSON unmarshals from JSON. null values become ""
func (s *CustomString) UnmarshalJSON(b []byte) error {
	return null.UnmarshalString(b, (*null.String)(s))
}

// Value returns the db value, null is returned for ""
func (s CustomString) Value() (driver.Value, error) {
	return null.String(s).Value()
}

// Scan scans from the db value. null values become ""
func (s *CustomString) Scan(value interface{}) error {
	return null.ScanString(value, (*null.String)(s))
}
```
