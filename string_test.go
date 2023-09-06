package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/nyaruka/null/v3"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value VARCHAR(255) NULL);`)

	tcs := []struct {
		value     null.String
		dbValue   driver.Value
		marshaled []byte
	}{
		{null.String("foo"), "foo", []byte(`"foo"`)},
		{null.NullString, nil, []byte(`null`)},
	}

	for _, tc := range tcs {
		mustExec(db, `DELETE FROM test`)

		dbValue, err := tc.value.Value()
		assert.NoError(t, err)
		assert.Equal(t, tc.dbValue, dbValue, "db value mismatch for %v", tc.value)

		// check writing the value to the database
		_, err = db.Exec(`INSERT INTO test(value) VALUES($1)`, tc.value)
		assert.NoError(t, err, "unexpected error writing %v", tc.value)

		rows, err := db.Query(`SELECT value FROM test;`)
		assert.NoError(t, err)

		var scanned null.String
		assert.True(t, rows.Next())
		err = rows.Scan(&scanned)
		assert.NoError(t, err)

		assert.Equal(t, tc.value, scanned, "scanned value mismatch for %v", tc.value)

		marshaled, err := json.Marshal(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.marshaled, marshaled, "marshaled mismatch for %v", tc.value)

		var unmarshaled null.String
		err = json.Unmarshal(marshaled, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, tc.value, unmarshaled, "unmarshaled mismatch for %v", tc.value)
	}
}

type CustomString string

const NullCustomString = CustomString("")

func (s *CustomString) Scan(value any) error         { return null.ScanString(value, s) }
func (s CustomString) Value() (driver.Value, error)  { return null.StringValue(s) }
func (s CustomString) MarshalJSON() ([]byte, error)  { return null.MarshalString(s) }
func (s *CustomString) UnmarshalJSON(b []byte) error { return null.UnmarshalString(b, s) }

func TestCustomString(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value VARCHAR(255) NULL);`)

	tcs := []struct {
		value     CustomString
		dbValue   driver.Value
		marshaled []byte
	}{
		{CustomString("foo"), "foo", []byte(`"foo"`)},
		{CustomString(""), nil, []byte(`null`)},
	}

	for _, tc := range tcs {
		mustExec(db, `DELETE FROM test`)

		dbValue, err := tc.value.Value()
		assert.NoError(t, err)
		assert.Equal(t, tc.dbValue, dbValue, "db value mismatch for %v", tc.value)

		// check writing the value to the database
		_, err = db.Exec(`INSERT INTO test(value) VALUES($1)`, tc.value)
		assert.NoError(t, err, "unexpected error writing %v", tc.value)

		rows, err := db.Query(`SELECT value FROM test;`)
		assert.NoError(t, err)

		var scanned CustomString
		assert.True(t, rows.Next())
		err = rows.Scan(&scanned)
		assert.NoError(t, err)

		assert.Equal(t, tc.value, scanned, "scanned value mismatch for %v", tc.value)

		marshaled, err := json.Marshal(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.marshaled, marshaled, "marshaled mismatch for %v", tc.value)

		var unmarshaled CustomString
		err = json.Unmarshal(marshaled, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, tc.value, unmarshaled, "unmarshaled mismatch for %v", tc.value)
	}
}
