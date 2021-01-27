package sqlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpsertBuilder(t *testing.T) {
	maker := Mysql{}.UpsertBuilder("table1", []string{"id", "status", "deleted"})

	sql, values, err := maker(
		[]interface{}{1, "succ", true},
	)
	assert.NoError(t, err)
	assert.Equal(t,
		"INSERT INTO `table1` (`id`,`status`,`deleted`) VALUES (?,?,?) ON DUPLICATE KEY UPDATE `id`=VALUES(`id`),`status`=VALUES(`status`),`deleted`=VALUES(`deleted`)",
		sql,
	)
	assert.Equal(t, []interface{}{1, "succ", true}, values)

	sql, values, err = maker(
		[]interface{}{1, "succ", true},
		[]interface{}{2, "fail", false},
	)
	assert.NoError(t, err)
	assert.Equal(t,
		"INSERT INTO `table1` (`id`,`status`,`deleted`) VALUES (?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `id`=VALUES(`id`),`status`=VALUES(`status`),`deleted`=VALUES(`deleted`)",
		sql,
	)
	assert.Equal(t, []interface{}{1, "succ", true, 2, "fail", false}, values)

	sql, values, err = maker(
		[]interface{}{3, "succ2", false},
		[]interface{}{4, "fail2", false},
		[]interface{}{5, "fail2", false},
	)
	assert.NoError(t, err)
	assert.Equal(t,
		"INSERT INTO `table1` (`id`,`status`,`deleted`) VALUES (?,?,?),(?,?,?),(?,?,?) ON DUPLICATE KEY UPDATE `id`=VALUES(`id`),`status`=VALUES(`status`),`deleted`=VALUES(`deleted`)",
		sql,
	)
	assert.Equal(t, []interface{}{3, "succ2", false, 4, "fail2", false, 5, "fail2", false}, values)
}
