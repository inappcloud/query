package query_test

import (
	"testing"

	"github.com/inappcloud/query"
)

func TestInsert(t *testing.T) {
	eq(t, "INSERT INTO teas DEFAULT VALUES", query.Insert("teas").String())

	eq(t, "INSERT INTO teas (name) VALUES(DEFAULT)", query.Insert("teas").Fields("name").String())
	eq(t, "INSERT INTO teas (name,color) VALUES(DEFAULT,DEFAULT)", query.Insert("teas").Fields("name,color").String())

	eq(t, "INSERT INTO teas (name) VALUES($1)", query.Insert("teas").Fields("name").Values("yame").String())
	eq(t, []interface{}{"yame"}, query.Insert("teas").Fields("name").Values("yame").Params())
	eq(t, []interface{}{"yame"}, query.Insert("teas").Fields("name").Values("yame", "other").Params())

	eq(t, "INSERT INTO teas (name,color) VALUES($1,$2)", query.Insert("teas").Fields("name,color").Values("yame", "green").String())
	eq(t, []interface{}{"yame", "green"}, query.Insert("teas").Fields("name,color").Values("yame", "green").Params())
	eq(t, "INSERT INTO teas (name,color) VALUES($1,DEFAULT)", query.Insert("teas").Fields("name,color").Values("yame").String())

	eq(t, "INSERT INTO teas DEFAULT VALUES RETURNING id", query.Insert("teas").Returning("id").String())
	eq(t, "INSERT INTO teas DEFAULT VALUES RETURNING id,name", query.Insert("teas").Returning("id,name").String())
}
