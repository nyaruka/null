package null_test

import (
	"encoding/json"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null/v2"
	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	db := getTestDB()

	_, err := db.Exec(`DROP TABLE IF EXISTS json_test; CREATE TABLE json_test(value jsonb null);`)
	assert.NoError(t, err)

	sp := func(s string) *string {
		return &s
	}

	tcs := []struct {
		Value null.JSON
		JSON  json.RawMessage
		DB    *string
	}{
		{null.JSON(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`), sp(`{"foo":"bar"}`)},
		{null.JSON(nil), json.RawMessage(`null`), nil},
		{null.JSON([]byte{}), json.RawMessage(`null`), nil},
	}

	for i, tc := range tcs {
		// first test marshalling and unmarshalling to JSON
		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, string(tc.JSON), string(b), "%d: marshalled json not equal", i)

		j := null.JSON("blah")
		err = json.Unmarshal(tc.JSON, &j)
		assert.NoError(t, err)
		assert.Equal(t, string(tc.Value), string(j), "%d: unmarshalled json not equal", i)

		// ok, now test writing and reading from DB
		_, err = db.Exec(`DELETE FROM json_test;`)
		assert.NoError(t, err)

		_, err = db.Exec(`INSERT INTO json_test(value) VALUES($1)`, tc.DB)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT value FROM json_test;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		j = null.JSON("blah")
		err = rows.Scan(&j)
		assert.NoError(t, err)

		if tc.Value == nil {
			assert.Nil(t, j, "%d: read db value should be null", i)
		} else {
			assert.Equal(t, string(tc.Value), strings.Replace(string(j), " ", "", -1), "%d: read db value should be equal", i)
		}

		_, err = db.Exec(`DELETE FROM json_test;`)
		assert.NoError(t, err)

		_, err = db.Exec(`INSERT INTO json_test(value) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err = db.Query(`SELECT value FROM json_test;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		var s *string
		err = rows.Scan(&s)
		assert.NoError(t, err)

		if tc.DB == nil {
			assert.Nil(t, s, "%d: written db value should be null", i)
		} else {
			assert.Equal(t, *tc.DB, strings.Replace(*s, " ", "", -1), "%d: written db value should be equal", i)
		}
	}
}
