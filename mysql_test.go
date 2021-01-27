package sqlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpsertMaker(t *testing.T) {
	maker := Mysql{}.UpsertMaker("table1", []string{"id", "status", "deleted"})
	sql, values, err := maker(
		[]interface{}{1, "succ", true},
		[]interface{}{2, "fail", false},
	)
	assert.NoError(t, err)
	assert.Equal(t,
		"INSERT INTO `table1` (`id`,`status`,`deleted`) VALUES (?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `id`=VALUES(`id`),`status`=VALUES(`status`),`deleted`=VALUES(`deleted`)",
		sql,
	)
	assert.Equal(t, []interface{}{1, "succ", true, 2, "fail", false}, values)
}
