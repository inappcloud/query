package query

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/inappcloud/query/where"
)

type updateQuery struct {
	table      string
	fields     []string
	values     []interface{}
	limit      int
	offset     int
	conditions *where.Condition
	returning  string
}

func (q *updateQuery) Fields(fields string) *updateQuery {
	if fields != "" {
		q.fields = strings.Split(fields, ",")
	}

	return q
}

func (q *updateQuery) Values(values ...interface{}) *updateQuery {
	if len(values) > 0 {
		q.values = values
	}

	return q
}

func (q *updateQuery) Limit(limit interface{}) *updateQuery {
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

func (q *updateQuery) Offset(offset interface{}) *updateQuery {
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

func (q *updateQuery) Where(c *where.Condition) *updateQuery {
	q.conditions = c
	return q
}

func (q *updateQuery) Returning(cols string) *updateQuery {
	if cols != "" {
		q.returning = cols
	}

	return q
}

func (q *updateQuery) String() string {
	buf := bytes.NewBufferString(fmt.Sprintf("UPDATE %s", q.table))

	if len(q.fields) > 0 {
		updates := make([]string, len(q.fields))
		for i, field := range q.fields {
			if i < len(q.values) {
				updates[i] = fmt.Sprintf("%s = ?", field)
			} else {
				updates[i] = fmt.Sprintf("%s = DEFAULT", field)
			}
		}

		buf.WriteString(" SET ")
		buf.WriteString(strings.Join(updates, ","))
	}

	if q.limit > 0 {
		buf.WriteString(fmt.Sprintf(" LIMIT %d", q.limit))
	}

	if q.offset > 0 {
		buf.WriteString(fmt.Sprintf(" OFFSET %d", q.offset))
	}

	if q.conditions != nil && len(q.conditions.String()) > 0 {
		buf.WriteString(fmt.Sprintf(" WHERE %s", q.conditions.String()))
	}

	if len(q.returning) > 0 {
		buf.WriteString(fmt.Sprintf(" RETURNING %s", q.returning))
	}

	return replacePlaceholders(buf.String())
}

func (q *updateQuery) Params() []interface{} {
	params := []interface{}{}

	if len(q.values) > 0 {
		params = append(params, q.values[:len(q.fields)]...)
	}

	if q.conditions != nil && len(q.conditions.Params()) > 0 {
		params = append(params, q.conditions.Params()...)
	}

	return params
}

func Update(table string) *updateQuery {
	return &updateQuery{table, []string{}, []interface{}{}, 0, 0, nil, ""}
}
