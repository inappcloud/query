package query

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/inappcloud/query/where"
)

type selectQuery struct {
	table      string
	fields     string
	limit      int
	offset     int
	conditions *where.Condition
}

func (q *selectQuery) Fields(fields []string) *selectQuery {
	if len(fields) > 0 {
		q.fields = strings.Join(fields, ",")
	}

	return q
}

func (q *selectQuery) Limit(limit interface{}) *selectQuery {
	switch val := limit.(type) {
	case int:
		q.limit = val
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			q.limit = i
		}
	}

	return q
}

func (q *selectQuery) Offset(offset interface{}) *selectQuery {
	switch val := offset.(type) {
	case int:
		q.offset = val
	case string:
		if i, err := strconv.Atoi(val); err == nil {
			q.offset = i
		}
	}

	return q
}

func (q *selectQuery) Where(c *where.Condition) *selectQuery {
	q.conditions = c
	return q
}

func (q *selectQuery) String() string {
	buf := bytes.NewBufferString(fmt.Sprintf("SELECT %s FROM %s", q.fields, q.table))

	if q.limit > 0 {
		buf.WriteString(fmt.Sprintf(" LIMIT %d", q.limit))
	}

	if q.offset > 0 {
		buf.WriteString(fmt.Sprintf(" OFFSET %d", q.offset))
	}

	if q.conditions != nil && len(q.conditions.String()) > 0 {
		buf.WriteString(fmt.Sprintf(" WHERE %s", q.conditions.String()))
	}

	return replacePlaceholders(buf.String())
}

func (q *selectQuery) Params() []interface{} {
	return q.conditions.Params()
}

func Select(table string) *selectQuery {
	return &selectQuery{table, "*", 0, 0, nil}
}
