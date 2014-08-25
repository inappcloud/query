package query

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/inappcloud/query/where"
)

type deleteQuery struct {
	table      string
	limit      int
	offset     int
	conditions *where.Condition
}

func (q *deleteQuery) Limit(limit interface{}) *deleteQuery {
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

func (q *deleteQuery) Offset(offset interface{}) *deleteQuery {
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

func (q *deleteQuery) Where(c *where.Condition) *deleteQuery {
	q.conditions = c
	return q
}

func (q *deleteQuery) String() string {
	buf := bytes.NewBufferString(fmt.Sprintf("DELETE FROM %s", q.table))

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

func (q *deleteQuery) Params() []interface{} {
	return q.conditions.Params()
}

func Delete(table string) *deleteQuery {
	return &deleteQuery{table, 0, 0, nil}
}
