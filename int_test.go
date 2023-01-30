package null_test

import (
	"database/sql/driver"
	"encoding/json"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nyaruka/null"
	"github.com/stretchr/testify/assert"
)

type CustomID int64

const NullCustomID = CustomID(0)

func (i *CustomID) Scan(value any) error {
	return null.ScanInt(value, i)
}

func (i CustomID) Value() (driver.Value, error) {
	return null.IntValue(i)
}

func (i *CustomID) UnmarshalJSON(b []byte) error {
	return null.UnmarshalInt(b, i)
}

func (i CustomID) MarshalJSON() ([]byte, error) {
	return null.MarshalInt(i)
}

type OtherCustom = null.Int64

func TestCustomInt64(t *testing.T) {
	db := getTestDB()

	_, err := db.Exec(`DROP TABLE IF EXISTS custom_id; CREATE TABLE custom_id(id integer null);`)
	assert.NoError(t, err)

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
		// {OtherCustom(10), "10", &ten}  // error, not the same type
	}

	for i, tc := range tcs {
		_, err = db.Exec(`DELETE FROM custom_id;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		id := CustomID(10)
		err = json.Unmarshal(b, &id)
		assert.NoError(t, err)
		assert.True(t, tc.Value == id, "%d: %s not equal to %s", i, tc.Value, id)
		assert.True(t, tc.Test == id, "%d: %s not equal to %s", i, tc.Test, id)

		_, err = db.Exec(`INSERT INTO custom_id(id) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT id FROM custom_id;`)
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

		rows, err = db.Query(`SELECT id FROM custom_id;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&id)
		assert.NoError(t, err)
		assert.True(t, tc.Value == id)
		assert.True(t, tc.Test == id)
	}
}

func TestInt64(t *testing.T) {
	db := getTestDB()

	_, err := db.Exec(`DROP TABLE IF EXISTS custom_id; CREATE TABLE custom_id(id integer null);`)
	assert.NoError(t, err)

	ten := int64(10)

	tcs := []struct {
		Value null.Int64
		JSON  string
		DB    *int64
	}{
		{null.Int64(10), "10", &ten},
		{null.Int64(0), "null", nil},
		{10, "10", &ten},
		// {OtherCustom(10), "10", &ten}  // error, not the same type
	}

	for i, tc := range tcs {
		_, err = db.Exec(`DELETE FROM custom_id;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		id := null.Int64(10)
		err = json.Unmarshal(b, &id)
		assert.NoError(t, err)
		assert.True(t, tc.Value == id)

		_, err = db.Exec(`INSERT INTO custom_id(id) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT id FROM custom_id;`)
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

		rows, err = db.Query(`SELECT id FROM custom_id;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&id)
		assert.NoError(t, err)
		assert.True(t, tc.Value == id)
	}
}
