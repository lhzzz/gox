package orm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOpenDB(t *testing.T) {
	db, err := OpenDB(
		MysqlOptions(MySQLOptions{
			User:     "dev",
			Password: "123456",
			Host:     "10.0.1.147",
			Port:     "3306",
			Database: "testDB",
		}),
		DebugMode(),
		SigularTable(),
		MaxLimit(100, 10, time.Hour),
	)
	assert.Nil(t, err)
	assert.EqualValues(t, db.Name(), "mysql")
}
