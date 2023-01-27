package null_test

import (
	"encoding/json"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	db := getTestDB()

	_, err := db.Exec(`DROP TABLE IF EXISTS map; CREATE TABLE map(value varchar(255) null);`)
	assert.NoError(t, err)

	sp := func(s string) *string {
		return &s
	}

	tcs := []struct {
		Value    null.Map
		JSON     string
		DB       *string
		Key      string
		KeyValue string
	}{
		{null.NewMap(map[string]interface{}{"foo": "bar"}), `{"foo":"bar"}`, sp(`{"foo": "bar"}`), "foo", "bar"},
		{null.NewMap(map[string]interface{}{}), "null", nil, "foo", ""},
		{null.NewMap(nil), "null", nil, "foo", ""},
		{null.NewMap(nil), "null", sp(""), "foo", ""},
	}

	for i, tc := range tcs {
		_, err = db.Exec(`DELETE FROM map;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, tc.JSON, string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		m := null.Map{}
		err = json.Unmarshal(b, &m)
		assert.NoError(t, err)
		assert.Equal(t, tc.Value.Map(), m.Map(), "%d: %s not equal to %s", i, tc.Value, m)
		assert.Equal(t, m.GetString(tc.Key, ""), tc.KeyValue)

		_, err = db.Exec(`INSERT INTO map(value) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT value FROM map;`)
		assert.NoError(t, err)

		m2 := null.Map{}
		assert.True(t, rows.Next())
		err = rows.Scan(&m2)
		assert.NoError(t, err)

		assert.Equal(t, tc.Value.Map(), m2.Map())
		assert.Equal(t, m2.GetString(tc.Key, ""), tc.KeyValue)

		_, err = db.Exec(`DELETE FROM map;`)
		assert.NoError(t, err)

		_, err = db.Exec(`INSERT INTO map(value) VALUES($1)`, tc.DB)
		assert.NoError(t, err)

		rows, err = db.Query(`SELECT value FROM map;`)
		assert.NoError(t, err)

		m2 = null.Map{}
		assert.True(t, rows.Next())
		err = rows.Scan(&m2)
		assert.NoError(t, err)

		assert.Equal(t, tc.Value.Map(), m2.Map())
		assert.Equal(t, m2.GetString(tc.Key, ""), tc.KeyValue)
	}
}
