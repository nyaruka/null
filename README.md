# null int and string values

This module provides (yet another) alternative in dealing with integers and strings which may be null in your JSON or 
database. There are various different approaches to this, namely the built in [golang SQL types](https://golang.org/pkg/database/sql/#NullInt64)
and the [guregu null module](https://github.com/guregu/null) which adds in JSON support. These are fine approaches but 
both suffer from you having to deal with a struct type for your ids instead of a more natural int64 or string. That is fine
in some cases but I prefer to use more primitive types as assignment is more straightforward, simple equality works, etc.

Sadly this requires a bit of boilerplate, so this package tries to make that a bit easier, proving helper methods for marshalling to/from
json and reading / scanning from a database value. Due to the Golang type system, using these requires some boilerplate, but I find it nicer
to have a bit more code for my ID types and less awkward code around the usage of ids rather than vice versa.

To define your own ID type which is nullable when written:

```golang
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

The process is essentially the same for nullable strings:

```golang
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

If you don't care about type safety within your types, you can use a go alias which simplifies things. This is more
useful for "stringy" things than it is for ids where you likely want to enforce type:

```golang
import "github.com/nyaruka/null"

// Note you lose type safety with aliases. FooID can be assigned to BarID and vice versa!
type FooID = null.Int
type BarID = null.Int

// Same here, any null.String can be assigned to CustomString
type CustomString = null.String
```
