package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null/v3"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	db := getTestDB()

	testMap := func() {
		tcs := []struct {
			value     null.JSON
			isNull    bool
			dbValue   driver.Value
			marshaled []byte
		}{
			{null.JSON(`{"foo": "bar"}`), false, []byte(`{"foo": "bar"}`), []byte(`{"foo": "bar"}`)},
			{null.JSON(`[1, 2, 3]`), false, []byte(`[1, 2, 3]`), []byte(`[1, 2, 3]`)},
			{null.JSON(`{}`), false, []byte(`{}`), []byte(`{}`)},
			{null.JSON(`null`), true, nil, []byte(`null`)},
			{null.JSON(``), true, nil, []byte(`null`)},
			{null.JSON(nil), true, nil, []byte(`null`)},
		}

		for _, tc := range tcs {
			mustExec(db, `DELETE FROM test`)

			assert.Equal(t, tc.isNull, tc.value.IsNull(), "isNull mismatch for %v", tc.value)

			dbValue, err := tc.value.Value()
			assert.NoError(t, err)
			assert.Equal(t, tc.dbValue, dbValue, "db value mismatch for %v", tc.value)

			// check writing the value to the database
			_, err = db.Exec(`INSERT INTO test(value) VALUES($1)`, tc.value)
			assert.NoError(t, err, "unexpected error writing %v", tc.value)

			rows, err := db.Query(`SELECT value FROM test;`)
			assert.NoError(t, err)

			scanned := null.JSON{}
			assert.True(t, rows.Next())
			err = rows.Scan(&scanned)
			assert.NoError(t, err)

			// we never return a nil JSON even if that's what we wrote
			expected := tc.value
			if len(expected) == 0 {
				expected = null.JSON(`null`)
			}

			assert.Equal(t, expected, scanned, "scanned value mismatch for %v", tc.value)

			marshaled, err := json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.JSONEq(t, string(tc.marshaled), string(marshaled), "marshaled mismatch for %v", tc.value)

			unmarshaled := null.JSON{}
			err = json.Unmarshal(marshaled, &unmarshaled)
			assert.NoError(t, err)
			assert.JSONEq(t, string(expected), string(unmarshaled), "unmarshaled mismatch for %v", tc.value)
		}
	}

	// test with TEXT column
	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value text null);`)
	testMap()

	// test with JSONB column
	mustExec(db, `DROP TABLE IF EXISTS test; CREATE TABLE test(value jsonb null);`)
	testMap()
}
