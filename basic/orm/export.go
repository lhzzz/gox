package orm

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func InitDB(dt DBType, opts ...Option) (*gorm.DB, error) {
	o := newMySQLOptions(dt)
	opts = append(opts, MysqlOptions(o))
	return OpenDB(opts...)
}

func newMySQLOptions(dbType DBType) MySQLOptions {
	return MySQLOptions{
		Master: MySQLDSNConfig{
			User:     viper.GetString(dbType.String() + ".UserName"),
			Password: viper.GetString(dbType.String() + ".Pwd"),
			Host:     viper.GetString(dbType.String() + ".HostName"),
			Port:     viper.GetString(dbType.String() + ".Port"),
			Database: viper.GetString(dbType.String() + ".DatabaseName"),
		},
	}
}
