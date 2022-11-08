package orm

import "time"

type DBType string

const (
	MySQL DBType = "MYSQL"
	TIDB  DBType = "TIDB"
)

func (d DBType) String() string {
	return string(d)
}

// MySQLOptions
type MySQLDSNConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type MySQLOptions struct {
	Master MySQLDSNConfig
	Slaves []MySQLDSNConfig
}

type dbOption struct {
	debugMode           bool
	sigularTable        bool
	dialect             DBType
	ignoreNotFoundError bool
	slowThreshold       time.Duration
	mysqlOptions        MySQLOptions
	maxOpenConn         int
	maxIdleConns        int
	maxLifetime         time.Duration
}

type Option func(do *dbOption)

var defaultDbOptions = dbOption{
	debugMode:           false,
	sigularTable:        false,
	dialect:             MySQL,
	slowThreshold:       time.Second,
	ignoreNotFoundError: true,
	maxOpenConn:         50,
	maxIdleConns:        10,
	maxLifetime:         time.Hour,
}

func (mdo *MySQLDSNConfig) DSN() string {
	return mdo.User + ":" + mdo.Password + "@tcp(" + mdo.Host + ":" + mdo.Port + ")/" + mdo.Database + "?charset=utf8mb4&parseTime=true&loc=Local"
}

func (do *MySQLOptions) MasterDSN() string {
	return do.Master.DSN()
}

func (do *MySQLOptions) SlaveDSNs() []string {
	res := make([]string, len(do.Slaves))
	for i := 0; i < len(do.Slaves); i++ {
		res[i] = do.Slaves[i].DSN()
	}
	return res
}

// MysqlOptions
func MysqlOptions(options MySQLOptions) Option {
	return func(o *dbOption) {
		o.dialect = MySQL
		o.mysqlOptions = options
	}
}

func TidbOptions(options MySQLOptions) Option {
	return func(o *dbOption) {
		o.dialect = TIDB
		o.mysqlOptions = options
	}
}

func DebugMode() Option {
	return func(do *dbOption) {
		do.debugMode = true
	}
}

func SigularTable() Option {
	return func(do *dbOption) {
		do.sigularTable = true
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

func MaxLimit(maxOpenConn, maxIdleConn int, maxLifetime time.Duration) Option {
	return func(do *dbOption) {
		do.maxOpenConn = maxOpenConn
		do.maxIdleConns = maxIdleConn
		do.maxLifetime = maxLifetime
	}
}
