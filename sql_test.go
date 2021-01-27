package sqlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	sql, values, err := Sql{}.In("id=? AND status in(?) OR deleted=?", 1, []string{"succ", "fail"}, true)
	assert.NoError(t, err)
	assert.Equal(t, "id=? AND status in(?, ?) OR deleted=?", sql)
	assert.Equal(t, []interface{}{1, "succ", "fail", true}, values)
}
