package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	"github.com/nyaruka/null/v2"
	"github.com/stretchr/testify/assert"
)

type CustomString string

const NullCustomString = CustomString("")

func (s *CustomString) Scan(value any) error         { return null.ScanString(value, s) }
func (s CustomString) Value() (driver.Value, error)  { return null.StringValue(s) }
func (s CustomString) MarshalJSON() ([]byte, error)  { return null.MarshalString(s) }
func (s *CustomString) UnmarshalJSON(b []byte) error { return null.UnmarshalString(b, s) }

func TestCustomString(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(string varchar(255) null);`)

	foo := "foo"

	tcs := []struct {
		Value CustomString
		JSON  string
		DB    *string
		Test  CustomString
	}{
		{CustomString(foo), `"foo"`, &foo, CustomString("foo")},
		{CustomString(""), "null", nil, NullCustomString},
		{NullCustomString, "null", nil, CustomString("")},
	}

	for _, tc := range tcs {
		mustExec(db, `DELETE FROM test;`)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%s not equal to %s", tc.JSON, string(b))

		str := CustomString("blah")
		err = json.Unmarshal(b, &str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str)
		assert.True(t, tc.Test == str)

		_, err = db.Exec(`INSERT INTO test(string) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT string FROM test;`)
		assert.NoError(t, err)

		var nullStr *string
		assert.True(t, rows.Next())
		err = rows.Scan(&nullStr)
		assert.NoError(t, err)

		if tc.DB == nil {
			assert.Nil(t, nullStr)
		} else {
			assert.True(t, *tc.DB == *nullStr)
		}

		rows, err = db.Query(`SELECT string FROM test;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str)
		assert.True(t, tc.Test == str)
	}
}

func TestString(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(string VARCHAR(255) NULL);`)

	foo := "foo"

	tcs := []struct {
		Value null.String
		JSON  string
		DB    *string
	}{
		{null.String("foo"), `"foo"`, &foo},
		{null.String(""), "null", nil},
		{null.NullString, "null", nil},
	}

	for i, tc := range tcs {
		mustExec(db, `DELETE FROM test;`)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		str := null.String("blah")
		err = json.Unmarshal(b, &str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str, "%d: %s not equal to %s", i, tc.Value, str)

		_, err = db.Exec(`INSERT INTO test(string) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT string FROM test;`)
		assert.NoError(t, err)

		var nullStr *string
		assert.True(t, rows.Next())
		err = rows.Scan(&nullStr)
		assert.NoError(t, err)

		if tc.DB == nil {
			assert.Nil(t, nullStr)
		} else {
			assert.True(t, *tc.DB == *nullStr)
		}

		rows, err = db.Query(`SELECT string FROM test;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str)
	}
}
