package where

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Condition struct {
	expression string
	params     []interface{}
}

func (c *Condition) valid() bool {
	return len(c.expression) > 0
}

func (c *Condition) String() string {
	return c.expression
}

func (c *Condition) Params() []interface{} {
	return c.params
}

func Eq(key string, value interface{}) *Condition {
	if value == nil {
		return &Condition{fmt.Sprintf("%s IS NULL", key), []interface{}{}}
	}

	return &Condition{fmt.Sprintf("%s = ?", key), []interface{}{value}}
}

func Ne(key string, value interface{}) *Condition {
	if value == nil {
		return &Condition{fmt.Sprintf("%s IS NOT NULL", key), []interface{}{}}
	}

	return &Condition{fmt.Sprintf("%s != ?", key), []interface{}{value}}
}

func Like(key string, value interface{}) *Condition {
	return &Condition{fmt.Sprintf("%s LIKE ?", key), []interface{}{value}}
}

func Gt(key string, value interface{}) *Condition {
	return &Condition{fmt.Sprintf("%s > ?", key), []interface{}{value}}
}

func Gte(key string, value interface{}) *Condition {
	return &Condition{fmt.Sprintf("%s >= ?", key), []interface{}{value}}
}

func Lt(key string, value interface{}) *Condition {
	return &Condition{fmt.Sprintf("%s < ?", key), []interface{}{value}}
}

func Lte(key string, value interface{}) *Condition {
	return &Condition{fmt.Sprintf("%s <= ?", key), []interface{}{value}}
}

func Between(key string, min interface{}, max interface{}) *Condition {
	return &Condition{fmt.Sprintf("%s BETWEEN ? AND ?", key), []interface{}{min, max}}
}

func In(key string, value []interface{}) *Condition {
	placeholders := strings.Repeat("?,", len(value))
	return &Condition{fmt.Sprintf("%s IN (%s)", key, placeholders[:len(placeholders)-1]), value}
}

func Nin(key string, value []interface{}) *Condition {
	placeholders := strings.Repeat("?,", len(value))
	return &Condition{fmt.Sprintf("%s NOT IN (%s)", key, placeholders[:len(placeholders)-1]), value}
}

func Or(conditions ...*Condition) *Condition {
	var expressions []string
	condition := &Condition{"", []interface{}{}}

	for _, c := range conditions {
		if c.valid() {
			expressions = append(expressions, c.expression)
			condition.params = append(condition.params, c.params...)
		}
	}

	if len(expressions) == 1 {
		condition.expression = expressions[0]
	}

	if len(expressions) > 1 {
		condition.expression = fmt.Sprintf("(%s)", strings.Join(expressions, " OR "))
	}

	return condition
}

func And(conditions ...*Condition) *Condition {
	var expressions []string
	condition := &Condition{"", []interface{}{}}

	for _, c := range conditions {
		if c.valid() {
			expressions = append(expressions, c.expression)
			condition.params = append(condition.params, c.params...)
		}
	}

	if len(expressions) == 1 {
		condition.expression = expressions[0]
	}

	if len(expressions) > 1 {
		condition.expression = fmt.Sprintf("(%s)", strings.Join(expressions, " AND "))
	}

	return condition
}

func Parse(v string) *Condition {
	if v == "" {
		return new(Condition)
	}

	buf := bytes.NewBufferString(v)
	where := make(map[string]interface{})
	json.NewDecoder(buf).Decode(&where)

	return addAndCondition(where)
}

func addAndCondition(condition map[string]interface{}) *Condition {
	var conds []*Condition

	for k, v := range condition {
		conds = append(conds, addCondition(k, v))
	}

	return And(conds...)
}

func addCondition(k string, v interface{}) *Condition {
	switch k {
	case "$or":
		var orConds []*Condition

		for _, condition := range v.([]interface{}) {
			orConds = append(orConds, addAndCondition(condition.(map[string]interface{})))
		}

		return Or(orConds...)
	case "$and":
		var andConds []*Condition

		for _, condition := range v.([]interface{}) {
			andConds = append(andConds, addAndCondition(condition.(map[string]interface{})))
		}

		return And(andConds...)
	default:
		switch typ := v.(type) {
		case map[string]interface{}:
			var conds []*Condition

			if val, ok := typ["$eq"]; ok {
				conds = append(conds, Eq(k, val))
			}

			if val, ok := typ["$ne"]; ok {
				conds = append(conds, Ne(k, val))
			}

			if val, ok := typ["$like"]; ok {
				conds = append(conds, Like(k, val))
			}

			if val, ok := typ["$gt"]; ok {
				conds = append(conds, Gt(k, val))
			}

			if val, ok := typ["$gte"]; ok {
				conds = append(conds, Gte(k, val))
			}

			if val, ok := typ["$lt"]; ok {
				conds = append(conds, Lt(k, val))
			}

			if val, ok := typ["$lte"]; ok {
				conds = append(conds, Lte(k, val))
			}

			// FIX: What if val is not an array or has less or more items than 2?
			// if val, ok := typ["$between"]; ok {
			// 	values := val.([]interface{})
			// 	conds = append(conds, Between(k, values[0], values[1]))
			// }

			return And(conds...)
		default:
			return Eq(k, v)
		}
	}

	return nil
}
