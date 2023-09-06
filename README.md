# null [![Build Status](https://github.com/nyaruka/null/workflows/CI/badge.svg)](https://github.com/nyaruka/null/actions?query=workflow%3ACI) [![codecov](https://codecov.io/gh/nyaruka/null/branch/main/graph/badge.svg)](https://codecov.io/gh/nyaruka/null) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/null)](https://goreportcard.com/report/github.com/nyaruka/null)

This module provides (yet another) alternative to dealing with null values in databases or JSON. Other approaches like 
the [Null types](https://golang.org/pkg/database/sql/#NullInt64) in the standard library use structs to ensure you can 
differentiate between zero values an null values. If that isn't a meaningful distinction in your app, then this module
might be a simpler approach for you because it uses primitive values and treats zero values as null values.

If you don't need to define your own types, you can use one of the following predefined types. If you scan a SQL `NULL` 
or unmarshal a JSON `null`, you will get the zero value. If you write the zero value to SQL you will get `NULL` and if
you marshal the zero value to JSON, you will get `null`.

|               | Zero Value
|---------------|-----------------
| `null.Int`    | `int(0)`        
| `null.Int64`  | `int64(0)`      
| `null.String` | `""`            
| `null.Map[V]`    | `map[string]V{}`         
| `null.JSON`   | `[]byte("null")`  

If you want to define a custom integer type, you need to define the following methods:

```go
import "github.com/nyaruka/null/v2"

type CustomID int64  // or int etc

const NullCustomID = CustomID(0)

func (i *CustomID) Scan(value any) error         { return null.ScanInt(value, i) }
func (i CustomID) Value() (driver.Value, error)  { return null.IntValue(i) }
func (i *CustomID) UnmarshalJSON(b []byte) error { return null.UnmarshalInt(b, i) }
func (i CustomID) MarshalJSON() ([]byte, error)  { return null.MarshalInt(i) }
```

And likewise for a custom string type:

```go
import "github.com/nyaruka/null/v2"

type CustomString string

const NullCustomString = CustomString("")

func (s *CustomString) Scan(value any) error         { return null.ScanString(value, s) }
func (s CustomString) Value() (driver.Value, error)  { return null.StringValue(s) }
func (s CustomString) MarshalJSON() ([]byte, error)  { return null.MarshalString(s) }
func (s *CustomString) UnmarshalJSON(b []byte) error { return null.UnmarshalString(b, s) }
```

If you want to create a type which can scan from `NULL`, but always writes as the zero value, just don't define the `Value` 
method. This can be useful when changing a database column to be non-NULL.
