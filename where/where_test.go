package where_test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/inappcloud/query/where"
)

func TestConditions(t *testing.T) {
	queryDataTable := []struct {
		condition  *where.Condition
		expression string
		params     []interface{}
	}{
		{where.Eq("name", "yame"), "name = ?", []interface{}{"yame"}},
		{where.Eq("name", nil), "name IS NULL", []interface{}{}},
		{where.Ne("name", "yame"), "name != ?", []interface{}{"yame"}},
		{where.Ne("name", nil), "name IS NOT NULL", []interface{}{}},
		{where.Like("name", "yame"), "name LIKE ?", []interface{}{"yame"}},
		{where.Gt("price", 10), "price > ?", []interface{}{10}},
		{where.Gte("price", 10), "price >= ?", []interface{}{10}},
		{where.Lt("price", 10), "price < ?", []interface{}{10}},
		{where.Lte("price", 10), "price <= ?", []interface{}{10}},
		{where.In("price", 10), "price IN (?)", []interface{}{10}},
		{where.In("price"), "", []interface{}{}},
		{where.In("price", 10, 11, 12), "price IN (?,?,?)", []interface{}{10, 11, 12}},
		{where.Nin("price", 10), "price NOT IN (?)", []interface{}{10}},
		{where.Nin("price"), "", []interface{}{}},
		{where.Nin("price", 10, 11, 12), "price NOT IN (?,?,?)", []interface{}{10, 11, 12}},
		{where.Between("price", 10, 20), "price BETWEEN ? AND ?", []interface{}{10, 20}},
		{where.Or(where.Eq("name", "yame"), where.Eq("color", "green")), "(name = ? OR color = ?)", []interface{}{"yame", "green"}},
		{where.And(where.Eq("name", "yame"), where.Eq("color", "green")), "(name = ? AND color = ?)", []interface{}{"yame", "green"}},
		{where.Or(where.Eq("name", "yame"), where.And(where.Eq("name", "phoenix"), where.Eq("color", "red"))), "(name = ? OR (name = ? AND color = ?))", []interface{}{"yame", "phoenix", "red"}},
		{where.And(new(where.Condition), where.Eq("name", "yame")), "name = ?", []interface{}{"yame"}},
		{where.Or(new(where.Condition), where.Eq("name", "yame")), "name = ?", []interface{}{"yame"}},
	}

	for _, test := range queryDataTable {
		eq(t, test.expression, test.condition.String())
		eq(t, test.params, test.condition.Params())
	}
}

func TestParse(t *testing.T) {
	eq(t, "", where.Parse("").String())
	eq(t, "", where.Parse("{}").String())
	eq(t, "", where.Parse("[]").String())

	whereDataTable := []struct {
		where     string
		expSql    string
		expParams []interface{}
	}{
		{
			`{"name":"yame"}`,
			"name = ?",
			[]interface{}{"yame"},
		},
		{
			`{"$or":[{"name":"yame"},{"color":"green"}]}`,
			"(name = ? OR color = ?)",
			[]interface{}{"yame", "green"},
		},
		{
			`{"$and":[{"name":"yame"},{"color":"green"}]}`,
			"(name = ? AND color = ?)",
			[]interface{}{"yame", "green"},
		},
		{
			`{"name":{"$eq":"yame"}}`,
			"name = ?",
			[]interface{}{"yame"},
		},
		{
			`{"name":{"$ne":"yame"}}`,
			"name != ?",
			[]interface{}{"yame"},
		},
		{
			`{"name":{"$eq":null}}`,
			"name IS NULL",
			[]interface{}{},
		},
		{
			`{"name":{"$ne":null}}`,
			"name IS NOT NULL",
			[]interface{}{},
		},
		{
			`{"name":{"$like":"ya%"}}`,
			"name LIKE ?",
			[]interface{}{"ya%"},
		},
		{
			`{"price":{"$gt":10}}`,
			"price > ?",
			[]interface{}{float64(10)},
		},
		{
			`{"price":{"$gte":10}}`,
			"price >= ?",
			[]interface{}{float64(10)},
		},
		{
			`{"price":{"$lt":10}}`,
			"price < ?",
			[]interface{}{float64(10)},
		},
		{
			`{"price":{"$lte":10}}`,
			"price <= ?",
			[]interface{}{float64(10)},
		},
		{
			`{"price":{"$between":[1,10]}}`,
			"price BETWEEN ? AND ?",
			[]interface{}{float64(1), float64(10)},
		},
		{
			`{"price":{"$between":[1,10,12]}}`,
			"",
			[]interface{}{},
		},
		{
			`{"price":{"$between":[1]}}`,
			"",
			[]interface{}{},
		},
		{
			`{"price":{"$between":[]}}`,
			"",
			[]interface{}{},
		},
		{
			`{"name":{"$in":[]}}`,
			"",
			[]interface{}{},
		},
		{
			`{"name":{"$in":["yame","uji"]}}`,
			"name IN (?,?)",
			[]interface{}{"yame", "uji"},
		},
		{
			`{"name":{"$nin":[]}}`,
			"",
			[]interface{}{},
		},
		{
			`{"name":{"$nin":["yame","uji"]}}`,
			"name NOT IN (?,?)",
			[]interface{}{"yame", "uji"},
		},
		{
			`{"$and":[{"name":{"$like":"ya%"}},{"color":"green"}]}`,
			"(name LIKE ? AND color = ?)",
			[]interface{}{"ya%", "green"},
		},
		{
			`{"$or":[{"name":{"$like":"ya%"}},{"$and":[{"name":"phoenix"},{"color":"red"}]}]}`,
			"(name LIKE ? OR (name = ? AND color = ?))",
			[]interface{}{"ya%", "phoenix", "red"},
		},
	}

	for _, test := range whereDataTable {
		c := where.Parse(test.where)
		eq(t, test.expSql, c.String())
		eq(t, test.expParams, c.Params())
	}

	c := where.Parse(`{"color":"green","name":"yame"}`)
	eqOr(t, "(color = ? AND name = ?)", "(name = ? AND color = ?)", c.String())
	eqOr(t, []interface{}{"green", "yame"}, []interface{}{"yame", "green"}, c.Params())
}

func eqOr(tb testing.TB, exp1, exp2, act interface{}) {
	if !reflect.DeepEqual(exp1, act) && !reflect.DeepEqual(exp2, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v or %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp1, exp2, act)
		tb.FailNow()
	}
}

func eq(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
