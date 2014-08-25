package query_test

import (
	"testing"

	"github.com/inappcloud/query"
	"github.com/inappcloud/query/where"
)

func TestDelete(t *testing.T) {
	eq(t, "DELETE FROM teas", query.Delete("teas").String())

	eq(t, "DELETE FROM teas LIMIT 1", query.Delete("teas").Limit(1).String())
	eq(t, "DELETE FROM teas LIMIT 1", query.Delete("teas").Limit("1").String())
	eq(t, "DELETE FROM teas", query.Delete("teas").Limit("not a number").String())
	eq(t, "DELETE FROM teas", query.Delete("teas").Limit(true).String())

	eq(t, "DELETE FROM teas OFFSET 1", query.Delete("teas").Offset(1).String())
	eq(t, "DELETE FROM teas OFFSET 1", query.Delete("teas").Offset("1").String())
	eq(t, "DELETE FROM teas", query.Delete("teas").Offset("not a number").String())
	eq(t, "DELETE FROM teas", query.Delete("teas").Offset(true).String())

	eq(t, "DELETE FROM teas", query.Delete("teas").Where(nil).String())
	eq(t, "DELETE FROM teas WHERE name = $1", query.Delete("teas").Where(where.Eq("name", "yame")).String())
	eq(t, []interface{}{"yame"}, query.Delete("teas").Where(where.Eq("name", "yame")).Params())
}
