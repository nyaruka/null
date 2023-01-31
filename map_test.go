package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	db := getTestDB()

	testMap := func() {
		tcs := []struct {
			value     null.Map
			dbValue   driver.Value
			marshaled []byte
		}{
			{null.Map{"foo": "bar"}, []byte(`{"foo":"bar"}`), []byte(`{"foo":"bar"}`)},
			{null.Map{}, nil, []byte(`null`)},
			{null.Map(nil), nil, []byte(`null`)},
		}

		for _, tc := range tcs {
			_, err := db.Exec(`DELETE FROM map;`)
			assert.NoError(t, err)

			dbValue, err := tc.value.Value()
			assert.NoError(t, err)
			assert.Equal(t, tc.dbValue, dbValue)

			// check writing the value to the database
			_, err = db.Exec(`INSERT INTO map(value) VALUES($1)`, tc.value)
			assert.NoError(t, err)

			rows, err := db.Query(`SELECT value FROM map;`)
			assert.NoError(t, err)

			scanned := null.Map{}
			assert.True(t, rows.Next())
			err = rows.Scan(&scanned)
			assert.NoError(t, err)

			// we never return a nil map even if that's what we wrote
			expected := tc.value
			if expected == nil {
				expected = null.Map{}
			}

			assert.Equal(t, expected, scanned)

			marshaled, err := json.Marshal(tc.value)
			assert.NoError(t, err)
			assert.Equal(t, tc.marshaled, marshaled)

			unmarshaled := null.Map{}
			err = json.Unmarshal(marshaled, &unmarshaled)
			assert.NoError(t, err)
			assert.Equal(t, expected, unmarshaled)
		}
	}

	mustExec(db, `DROP TABLE IF EXISTS map; CREATE TABLE map(value text null);`)

	testMap()

	mustExec(db, `DROP TABLE IF EXISTS map; CREATE TABLE map(value jsonb null);`)

	testMap()
}
