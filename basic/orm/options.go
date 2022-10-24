package orm

import "time"

type DBType string

const (
	MySQL DBType = "MYSQL"
)

// MySQLOptions
type MySQLOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type dbOption struct {
	debugMode           bool
	sigularTable        bool
	dialect             DBType
	ignoreNotFoundError bool
	slowThreshold       time.Duration
	mysqlOptions        MySQLOptions
}

type Option func(do *dbOption)

var defaultDbOptions = dbOption{
	debugMode:           false,
	sigularTable:        false,
	dialect:             MySQL,
	slowThreshold:       time.Second,
	ignoreNotFoundError: true,
}

func (do *dbOption) MySQLDSN() string {
	return do.mysqlOptions.User + ":" + do.mysqlOptions.Password + "@tcp(" + do.mysqlOptions.Host + ":" + do.mysqlOptions.Port + ")/" + do.mysqlOptions.Database + "?charset=utf8mb4&parseTime=true&loc=Local"
}

// MysqlOptions
func MysqlOptions(options MySQLOptions) Option {
	return func(o *dbOption) {
		o.dialect = MySQL
		o.mysqlOptions = options
	}
}

func SetDebug(debug bool) Option {
	return func(do *dbOption) {
		do.debugMode = debug
	}
}

func SetSigularTable(enable bool) Option {
	return func(do *dbOption) {
		do.sigularTable = enable
	}
}

func EnableRecordNotFoundError() Option {
	return func(do *dbOption) {
		do.ignoreNotFoundError = false
	}
}

func SlowThreshold(d time.Duration) Option {
	return func(do *dbOption) {
		do.slowThreshold = d
	}
}
