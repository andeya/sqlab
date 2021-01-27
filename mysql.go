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

func (m *Mysql) UpsertMaker(
	tableName string, fields []string,
) func(valueGroups ...[]interface{}) (sql string, values []interface{}, err error) {
	if len(fields) == 0 {
		return func(valueGroups ...[]interface{}) (sql string, values []interface{}, err error) {
			return "", nil, errors.New("no fields")
		}
	}
	if !strings.HasPrefix(tableName, "`") {
		tableName = "`" + tableName + "`"
	}
	updateFields := make([]string, len(fields))
	for i, field := range fields {
		if !strings.HasPrefix(field, "`") {
			field = "`" + field + "`"
			fields[i] = field
		}
		updateFields[i] = field + "=VALUES(" + field + ")"
	}
	var sqlTpl = fmt.Sprintf("INSERT INTO %s (%s) VALUES %%s ON DUPLICATE KEY UPDATE %s",
		tableName,
		strings.Join(fields, ","),
		strings.Join(updateFields, ","),
	)
	return func(valueGroups ...[]interface{}) (sql string, values []interface{}, err error) {
		if len(valueGroups) == 0 {
			return "", nil, errors.New("no data to update")
		}
		values = make([]interface{}, 0, len(fields)*len(valueGroups))
		for _, vars := range valueGroups {
			if len(vars) != len(fields) {
				return "", nil, fmt.Errorf(
					"number of values and fields is not equal: %d!=%d",
					len(vars), len(fields),
				)
			}
			values = append(values, vars...)
		}
		sql = fmt.Sprintf(sqlTpl,
			strings.Repeat(",("+strings.Repeat(",?", len(fields))[1:]+")", len(valueGroups))[1:],
		)
		return sql, values, nil
	}
}
