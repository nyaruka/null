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

func TestJSONNullIsNotShared(t *testing.T) {
	// scanning a NULL must not hand out a value which aliases package state or any other value
	var v1, v2 null.JSON
	assert.NoError(t, v1.Scan(nil))
	assert.NoError(t, v2.Scan(nil))

	// unmarshaling a short value into one must not write through to the other or to null.NullJSON
	assert.NoError(t, json.Unmarshal([]byte(`12`), &v1))

	assert.Equal(t, null.JSON(`12`), v1)
	assert.Equal(t, null.JSON(`null`), v2)
	assert.True(t, v2.IsNull())
	assert.Equal(t, null.JSON(`null`), null.NullJSON)

	// and a NULL scanned afterwards is still null
	var v3 null.JSON
	assert.NoError(t, v3.Scan(nil))
	assert.Equal(t, null.JSON(`null`), v3)
	assert.True(t, v3.IsNull())

	// a null value still writes as SQL NULL rather than the literal bytes
	dbValue, err := null.JSON(`null`).Value()
	assert.NoError(t, err)
	assert.Nil(t, dbValue)

	// same for a value scanned from empty bytes
	var v4 null.JSON
	assert.NoError(t, v4.Scan([]byte{}))
	assert.NoError(t, json.Unmarshal([]byte(`[]`), &v4))
	assert.Equal(t, null.JSON(`null`), null.NullJSON)
}

func TestJSONIsNullDoesntDependOnNullJSON(t *testing.T) {
	orig := null.NullJSON
	defer func() { null.NullJSON = orig }()

	null.NullJSON = null.JSON(`{"not":"null"}`)

	assert.True(t, null.JSON(`null`).IsNull())
	assert.True(t, null.JSON(``).IsNull())
	assert.True(t, null.JSON(nil).IsNull())
	assert.False(t, null.JSON(`{"not":"null"}`).IsNull())
}

func TestJSONMarshalDoesntAliasValue(t *testing.T) {
	v := null.JSON(`{"foo":"bar"}`)

	marshaled, err := v.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"foo":"bar"}`), marshaled)

	marshaled[2] = 'X'

	assert.Equal(t, null.JSON(`{"foo":"bar"}`), v)
}

func TestJSONScanDoesntAliasDriverBuffer(t *testing.T) {
	// the driver is allowed to reuse the buffer it gives us for subsequent scans
	buf := []byte(`{"foo":"bar"}`)

	var v null.JSON
	assert.NoError(t, v.Scan(buf))

	copy(buf, []byte(`{"zzz":"zzz"}`))

	assert.Equal(t, null.JSON(`{"foo":"bar"}`), v)
}
