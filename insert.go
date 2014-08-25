package query

import (
	"bytes"
	"fmt"
	"strings"
)

type insertQuery struct {
	table     string
	fields    []string
	values    []interface{}
	returning string
}

func (q *insertQuery) Fields(fields string) *insertQuery {
	if fields != "" {
		q.fields = strings.Split(fields, ",")
	}

	return q
}

func (q *insertQuery) Values(values ...interface{}) *insertQuery {
	if len(values) > 0 {
		q.values = values
	}

	return q
}

func (q *insertQuery) Returning(cols string) *insertQuery {
	if cols != "" {
		q.returning = cols
	}

	return q
}

func (q *insertQuery) String() string {
	buf := bytes.NewBufferString(fmt.Sprintf("INSERT INTO %s", q.table))

	if len(q.fields) > 0 {

		values := make([]string, len(q.fields))
		for i := range q.fields {
			if i < len(q.values) {
				values[i] = "?"
			} else {
				values[i] = "DEFAULT"
			}
		}

		buf.WriteString(fmt.Sprintf(" (%s) VALUES(%s)", strings.Join(q.fields, ","), strings.Join(values, ",")))

	} else {
		buf.WriteString(" DEFAULT VALUES")
	}

	if len(q.returning) > 0 {
		buf.WriteString(fmt.Sprintf(" RETURNING %s", q.returning))
	}

	return replacePlaceholders(buf.String())
}

func (q *insertQuery) Params() []interface{} {
	params := []interface{}{}

	if len(q.values) > 0 {
		params = append(params, q.values[:len(q.fields)]...)
	}

	return params
}

func Insert(table string) *insertQuery {
	return &insertQuery{table, []string{}, []interface{}{}, ""}
}
