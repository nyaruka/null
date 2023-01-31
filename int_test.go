package null_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null/v2"
	"github.com/stretchr/testify/assert"
)

func TestInt(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value VARCHAR(255) NULL);`)

	tcs := []struct {
		value     null.Int
		dbValue   driver.Value
		marshaled []byte
	}{
		{null.Int(123), int64(123), []byte(`123`)},
		{null.NullInt, nil, []byte(`null`)},
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

		var scanned null.Int
		assert.True(t, rows.Next())
		err = rows.Scan(&scanned)
		assert.NoError(t, err)

		assert.Equal(t, tc.value, scanned, "scanned value mismatch for %v", tc.value)

		marshaled, err := json.Marshal(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.marshaled, marshaled, "marshaled mismatch for %v", tc.value)

		var unmarshaled null.Int
		err = json.Unmarshal(marshaled, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, tc.value, unmarshaled, "unmarshaled mismatch for %v", tc.value)
	}
}

func TestInt64(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value VARCHAR(255) NULL);`)

	tcs := []struct {
		value     null.Int64
		dbValue   driver.Value
		marshaled []byte
	}{
		{null.Int64(123), int64(123), []byte(`123`)},
		{null.NullInt64, nil, []byte(`null`)},
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

		var scanned null.Int64
		assert.True(t, rows.Next())
		err = rows.Scan(&scanned)
		assert.NoError(t, err)

		assert.Equal(t, tc.value, scanned, "scanned value mismatch for %v", tc.value)

		marshaled, err := json.Marshal(tc.value)
		assert.NoError(t, err)
		assert.Equal(t, tc.marshaled, marshaled, "marshaled mismatch for %v", tc.value)

		var unmarshaled null.Int64
		err = json.Unmarshal(marshaled, &unmarshaled)
		assert.NoError(t, err)
		assert.Equal(t, tc.value, unmarshaled, "unmarshaled mismatch for %v", tc.value)
	}
}

type CustomID int64

const NullCustomID = CustomID(0)

func (i *CustomID) Scan(value any) error         { return null.ScanInt(value, i) }
func (i CustomID) Value() (driver.Value, error)  { return null.IntValue(i) }
func (i *CustomID) UnmarshalJSON(b []byte) error { return null.UnmarshalInt(b, i) }
func (i CustomID) MarshalJSON() ([]byte, error)  { return null.MarshalInt(i) }

func TestCustomInt64(t *testing.T) {
	db := getTestDB()

	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(id INT NULL);`)

	ten := int64(10)

	tcs := []struct {
		Value CustomID
		JSON  string
		DB    *int64
		Test  CustomID
	}{
		{CustomID(10), "10", &ten, CustomID(10)},
		{CustomID(0), "null", nil, NullCustomID},
		{10, "10", &ten, CustomID(10)},
		{NullCustomID, "null", nil, CustomID(0)},
	}

	for i, tc := range tcs {
		mustExec(db, `DELETE FROM test`)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		id := CustomID(10)
		err = json.Unmarshal(b, &id)
		assert.NoError(t, err)
		assert.True(t, tc.Value == id, "%d: %s not equal to %s", i, tc.Value, id)
		assert.True(t, tc.Test == id, "%d: %s not equal to %s", i, tc.Test, id)

		_, err = db.Exec(`INSERT INTO test(id) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT id FROM test;`)
		assert.NoError(t, err)

		var intID *int64
		assert.True(t, rows.Next())
		err = rows.Scan(&intID)
		assert.NoError(t, err)

		if tc.DB == nil {
			assert.Nil(t, intID)
		} else {
			assert.True(t, *tc.DB == *intID)
		}

		rows, err = db.Query(`SELECT id FROM test;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&id)
		assert.NoError(t, err)
		assert.True(t, tc.Value == id)
		assert.True(t, tc.Test == id)
	}
}

func getTestDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://nyaruka:nyaruka@localhost/null_test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return db
}

func mustExec(db *sql.DB, q string) sql.Result {
	res, err := db.Exec(q)
	if err != nil {
		panic(err)
	}
	return res
}
