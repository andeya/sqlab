package sqlab

import (
	"errors"
	"fmt"
	"strings"
)

// Mysql mysql tools
type Mysql struct {
	Sql
}

// UpsertBuilder build INSERT INTO ... ON DUPLICATE KEY UPDATE ...
func (m Mysql) UpsertBuilder(
	tableName string, fields []string,
) func(valueGroups ...[]interface{}) (sql string, values []interface{}, err error) {
	if len(tableName) == 0 || len(fields) == 0 {
		return func(valueGroups ...[]interface{}) (sql string, values []interface{}, err error) {
			return "", nil, errors.New("no table name or no fields")
		}
	}
	buf := strings.Builder{}
	buf.Grow(20 + len(tableName) + (len(fields) * 5)) // experience value
	buf.WriteString("INSERT INTO ")
	if tableName[0] != '`' {
		buf.WriteByte('`')
		buf.WriteString(tableName)
		buf.WriteByte('`')
	} else {
		buf.WriteString(tableName)
	}
	buf.WriteString(" (")
	buf2 := strings.Builder{}
	buf2.Grow(30 + (len(fields)*6)*3) // experience value
	buf2.WriteString(" ON DUPLICATE KEY UPDATE ")
	for i, field := range fields {
		if !strings.HasPrefix(field, "`") {
			buf.WriteByte('`')
			buf.WriteString(field)
			buf.WriteByte('`')
			buf2.WriteByte('`')
			buf2.WriteString(field)
			buf2.WriteString("`=VALUES(`")
			buf2.WriteString(field)
			buf2.WriteString("`)")
		} else {
			buf.WriteString(field)
			buf2.WriteString(field)
			buf2.WriteString("=VALUES(")
			buf2.WriteString(field)
			buf2.WriteString(")")
		}
		if i < len(fields)-1 {
			buf.WriteByte(',')
			buf2.WriteByte(',')
		}
	}
	buf.WriteString(") VALUES ")
	segment1 := buf.String()
	groupMark := "(" + strings.Repeat(",?", len(fields))[1:] + ")"
	segment2 := buf2.String()
	return func(valueGroups ...[]interface{}) (sql string, values []interface{}, err error) {
		if len(valueGroups) == 0 {
			return "", nil, errors.New("no data to update")
		}
		values = make([]interface{}, 0, len(fields)*len(valueGroups))
		bufOne := strings.Builder{}
		bufOne.Grow(len(segment1) + len(segment2) + (len(valueGroups) * (len(groupMark) + 1)))
		bufOne.WriteString(segment1)
		for i, vars := range valueGroups {
			if len(vars) != len(fields) {
				return "", nil, fmt.Errorf(
					"number of values and fields is not equal: %d!=%d",
					len(vars), len(fields),
				)
			}
			values = append(values, vars...)
			bufOne.WriteString(groupMark)
			if i < len(valueGroups)-1 {
				bufOne.WriteByte(',')
			}
		}
		bufOne.WriteString(segment2)
		sql = bufOne.String()
		return sql, values, nil
	}
}
