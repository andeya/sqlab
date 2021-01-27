package sqlab

import (
	"github.com/henrylee2cn/sqlab/internal/sqlx"
)

// Sql sql common tools
type Sql struct {
}

// In expands slice values in args, returning the modified query string
// and a new arg list that can be executed by a database. The `query` should
// use the `?` bindVar.  The return value uses the `?` bindVar.
func (s *Sql) In(query string, args ...interface{}) (string, []interface{}, error) {
	return sqlx.In(query, args)
}
