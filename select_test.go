package query_test

import (
	"testing"

	"github.com/inappcloud/query"
	"github.com/inappcloud/query/where"
)

func TestSelect(t *testing.T) {
	eq(t, "SELECT * FROM teas", query.Select("teas").String())

	eq(t, "SELECT name FROM teas", query.Select("teas").Fields([]string{"name"}).String())
	eq(t, "SELECT * FROM teas", query.Select("teas").Fields([]string{}).String())

	eq(t, "SELECT * FROM teas LIMIT 1", query.Select("teas").Limit(1).String())
	eq(t, "SELECT * FROM teas LIMIT 1", query.Select("teas").Limit("1").String())
	eq(t, "SELECT * FROM teas", query.Select("teas").Limit("not a number").String())
	eq(t, "SELECT * FROM teas", query.Select("teas").Limit(true).String())

	eq(t, "SELECT * FROM teas OFFSET 1", query.Select("teas").Offset(1).String())
	eq(t, "SELECT * FROM teas OFFSET 1", query.Select("teas").Offset("1").String())
	eq(t, "SELECT * FROM teas", query.Select("teas").Offset("not a number").String())
	eq(t, "SELECT * FROM teas", query.Select("teas").Offset(true).String())

	eq(t, "SELECT * FROM teas", query.Select("teas").Where(nil).String())
	eq(t, "SELECT * FROM teas WHERE name = $1", query.Select("teas").Where(where.Eq("name", "yame")).String())
	eq(t, []interface{}{"yame"}, query.Select("teas").Where(where.Eq("name", "yame")).Params())

	eq(t, "SELECT * FROM teas ORDER BY name DESC, id ASC", query.Select("teas").Sort(map[string]int{"name": -1, "id": 1}).String())

	eq(t, "SELECT COUNT(*) FROM teas", query.Select("teas").Count())
}
