package query_test

import (
	"testing"

	"github.com/inappcloud/query"
	"github.com/inappcloud/query/where"
)

func TestUpdate(t *testing.T) {
	eq(t, "UPDATE teas", query.Update("teas").String())

	eq(t, "UPDATE teas SET name = DEFAULT", query.Update("teas").Fields("name").String())
	eq(t, "UPDATE teas SET name = DEFAULT,color = DEFAULT", query.Update("teas").Fields("name,color").String())

	eq(t, "UPDATE teas SET name = $1", query.Update("teas").Fields("name").Values("yame").String())
	eq(t, []interface{}{"yame"}, query.Update("teas").Fields("name").Values("yame").Params())
	eq(t, []interface{}{"yame"}, query.Update("teas").Fields("name").Values("yame", "other").Params())

	eq(t, "UPDATE teas SET name = $1,color = $2", query.Update("teas").Fields("name,color").Values("yame", "green").String())
	eq(t, []interface{}{"yame", "green"}, query.Update("teas").Fields("name,color").Values("yame", "green").Params())
	eq(t, "UPDATE teas SET name = $1,color = DEFAULT", query.Update("teas").Fields("name,color").Values("yame").String())

	eq(t, "UPDATE teas LIMIT 1", query.Update("teas").Limit(1).String())
	eq(t, "UPDATE teas LIMIT 1", query.Update("teas").Limit("1").String())
	eq(t, "UPDATE teas", query.Update("teas").Limit("not a number").String())
	eq(t, "UPDATE teas", query.Update("teas").Limit(true).String())

	eq(t, "UPDATE teas OFFSET 1", query.Update("teas").Offset(1).String())
	eq(t, "UPDATE teas OFFSET 1", query.Update("teas").Offset("1").String())
	eq(t, "UPDATE teas", query.Update("teas").Offset("not a number").String())
	eq(t, "UPDATE teas", query.Update("teas").Offset(true).String())

	eq(t, "UPDATE teas", query.Update("teas").Where(nil).String())
	eq(t, "UPDATE teas WHERE name = $1", query.Update("teas").Where(where.Eq("name", "yame")).String())
	eq(t, []interface{}{"yame"}, query.Update("teas").Where(where.Eq("name", "yame")).Params())

	eq(t, "UPDATE teas RETURNING id", query.Update("teas").Returning("id").String())
	eq(t, "UPDATE teas RETURNING id,name", query.Update("teas").Returning("id,name").String())
}
