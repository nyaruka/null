package null

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
)

type CustomID Int

func (i CustomID) MarshalJSON() ([]byte, error) {
	return Int(i).MarshalJSON()
}

func (i *CustomID) UnmarshalJSON(b []byte) error {
	return UnmarshalInt(b, (*Int)(i))
}

func (i CustomID) Value() (driver.Value, error) {
	return Int(i).Value()
}

func (i *CustomID) Scan(value interface{}) error {
	return ScanInt(value, (*Int)(i))
}

type OtherCustom = Int

const NullCustomID = CustomID(0)

func TestCustomInt(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost/null_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`DROP TABLE IF EXISTS custom_id; CREATE TABLE custom_id(id integer null);`)
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

func TestInt(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost/null_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`DROP TABLE IF EXISTS custom_id; CREATE TABLE custom_id(id integer null);`)
	assert.NoError(t, err)

	ten := int64(10)

	tcs := []struct {
		Value Int
		JSON  string
		DB    *int64
	}{
		{Int(10), "10", &ten},
		{Int(0), "null", nil},
		{10, "10", &ten},
		// {OtherCustom(10), "10", &ten}  // error, not the same type
	}

	for i, tc := range tcs {
		_, err = db.Exec(`DELETE FROM custom_id;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		id := Int(10)
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

type CustomString String

func (s CustomString) MarshalJSON() ([]byte, error) {
	return String(s).MarshalJSON()
}

func (s *CustomString) UnmarshalJSON(b []byte) error {
	return UnmarshalString(b, (*String)(s))
}

func (s CustomString) Value() (driver.Value, error) {
	return String(s).Value()
}

func (s *CustomString) Scan(value interface{}) error {
	return ScanString(value, (*String)(s))
}

const NullCustomString = CustomString("")

func TestCustomString(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost/null_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`DROP TABLE IF EXISTS custom_string; CREATE TABLE custom_string(string varchar(255) null);`)
	assert.NoError(t, err)

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
		_, err = db.Exec(`DELETE FROM custom_string;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%s not equal to %s", tc.JSON, string(b))

		str := CustomString("blah")
		err = json.Unmarshal(b, &str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str)
		assert.True(t, tc.Test == str)

		_, err = db.Exec(`INSERT INTO custom_string(string) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT string FROM custom_string;`)
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

		rows, err = db.Query(`SELECT string FROM custom_string;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str)
		assert.True(t, tc.Test == str)
	}
}

func TestString(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost/null_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`DROP TABLE IF EXISTS custom_string; CREATE TABLE custom_string(string varchar(255) null);`)
	assert.NoError(t, err)

	foo := "foo"

	tcs := []struct {
		Value String
		JSON  string
		DB    *string
	}{
		{String("foo"), `"foo"`, &foo},
		{String(""), "null", nil},
		{NullString, "null", nil},
	}

	for i, tc := range tcs {
		_, err = db.Exec(`DELETE FROM custom_string;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.True(t, tc.JSON == string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		str := String("blah")
		err = json.Unmarshal(b, &str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str, "%d: %s not equal to %s", i, tc.Value, str)

		_, err = db.Exec(`INSERT INTO custom_string(string) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT string FROM custom_string;`)
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

		rows, err = db.Query(`SELECT string FROM custom_string;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		err = rows.Scan(&str)
		assert.NoError(t, err)
		assert.True(t, tc.Value == str)
	}
}

func TestMap(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost/null_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`DROP TABLE IF EXISTS map; CREATE TABLE map(value varchar(255) null);`)
	assert.NoError(t, err)

	sp := func(s string) *string {
		return &s
	}

	tcs := []struct {
		Value    Map
		JSON     string
		DB       *string
		Key      string
		KeyValue string
	}{
		{NewMap(map[string]interface{}{"foo": "bar"}), `{"foo":"bar"}`, sp(`{"foo": "bar"}`), "foo", "bar"},
		{NewMap(map[string]interface{}{}), "null", nil, "foo", ""},
		{NewMap(nil), "null", nil, "foo", ""},
		{NewMap(nil), "null", sp(""), "foo", ""},
	}

	for i, tc := range tcs {
		_, err = db.Exec(`DELETE FROM map;`)
		assert.NoError(t, err)

		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, tc.JSON, string(b), "%d: %s not equal to %s", i, tc.JSON, string(b))

		m := Map{}
		err = json.Unmarshal(b, &m)
		assert.NoError(t, err)
		assert.Equal(t, tc.Value.Map(), m.Map(), "%d: %s not equal to %s", i, tc.Value, m)
		assert.Equal(t, m.GetString(tc.Key, ""), tc.KeyValue)

		_, err = db.Exec(`INSERT INTO map(value) VALUES($1)`, tc.Value)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT value FROM map;`)
		assert.NoError(t, err)

		m2 := Map{}
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

		m2 = Map{}
		assert.True(t, rows.Next())
		err = rows.Scan(&m2)
		assert.NoError(t, err)

		assert.Equal(t, tc.Value.Map(), m2.Map())
		assert.Equal(t, m2.GetString(tc.Key, ""), tc.KeyValue)
	}
}

func TestJSON(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://localhost/null_test?sslmode=disable")
	assert.NoError(t, err)

	_, err = db.Exec(`DROP TABLE IF EXISTS json_test; CREATE TABLE json_test(value jsonb null);`)
	assert.NoError(t, err)

	sp := func(s string) *string {
		return &s
	}

	tcs := []struct {
		Value     JSON
		WriteJSON json.RawMessage
		ReadJSON  json.RawMessage
		WriteDB   *string
		ReadDB    *string
	}{
		{JSON(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`), json.RawMessage(`{"foo":"bar"}`), sp(`{"foo":"bar"}`), sp(`{"foo":"bar"}`)},
		{JSON(nil), json.RawMessage(`null`), json.RawMessage(`null`), nil, nil},
		{JSON(`{}`), json.RawMessage(`{}`), json.RawMessage(`{}`), sp(`{}`), sp(`{}`)},
		{JSON(nil), json.RawMessage(`null`), json.RawMessage(`null`), nil, sp(`null`)},
		{JSON([]byte{}), json.RawMessage(`null`), json.RawMessage(`null`), nil, nil},
	}

	for i, tc := range tcs {
		// first test marshalling and unmarshalling to JSON
		b, err := json.Marshal(tc.Value)
		assert.NoError(t, err)
		assert.Equal(t, string(tc.WriteJSON), string(b), "%d: marshalled json not equal", i)

		j := JSON("blah")
		err = json.Unmarshal(tc.ReadJSON, &j)
		assert.NoError(t, err)
		assert.Equal(t, string(tc.Value), string(j), "%d: unmarshalled json not equal", i)

		// ok, now test writing and reading from DB
		_, err = db.Exec(`DELETE FROM json_test;`)
		assert.NoError(t, err)

		_, err = db.Exec(`INSERT INTO json_test(value) VALUES($1)`, tc.ReadDB)
		assert.NoError(t, err)

		rows, err := db.Query(`SELECT value FROM json_test;`)
		assert.NoError(t, err)

		assert.True(t, rows.Next())
		j = JSON("blah")
		err = rows.Scan(&j)
		assert.NoError(t, err)

		if tc.Value == nil {
			assert.Nil(t, j)
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

		if tc.WriteDB == nil {
			assert.Nil(t, s)
		} else {
			assert.Equal(t, *tc.WriteDB, strings.Replace(*s, " ", "", -1), "%d: written db value should be equal", i)
		}
	}
}
