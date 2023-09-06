package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null/v2"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	db := getTestDB()

	testMap := func() {
		tcs := []struct {
			value     null.Map[string]
			dbValue   driver.Value
			marshaled []byte
		}{
			{null.Map[string]{"foo": "bar"}, []byte(`{"foo":"bar"}`), []byte(`{"foo":"bar"}`)},
			{null.Map[string]{}, nil, []byte(`null`)},
			{null.Map[string](nil), nil, []byte(`null`)},
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

			scanned := null.Map[string]{}
			assert.True(t, rows.Next())
			err = rows.Scan(&scanned)
			assert.NoError(t, err)

			// we never return a nil map even if that's what we wrote
			expected := tc.value
			if expected == nil {
				expected = null.Map[string]{}
			}

			assert.Equal(t, expected, scanned, "scanned value mismatch for %v", tc.value)

			marshaled, err := json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.marshaled, marshaled, "marshaled mismatch for %v", tc.value)

			unmarshaled := null.Map[string]{}
			err = json.Unmarshal(marshaled, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, expected, unmarshaled, "unmarshaled mismatch for %v", tc.value)
		}
	}

	// test with TEXT column
	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value text null);`)
	testMap()

	// test with JSONB column
	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value jsonb null);`)
	testMap()
}
