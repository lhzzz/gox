package orm

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	gormopentracing "gorm.io/plugin/opentracing"
)

func OpenDB(opt ...Option) (*gorm.DB, error) {
	opts := defaultDbOptions
	for _, o := range opt {
		o(&opts)
	}

	lvl := logger.Warn
	if opts.debugMode {
		lvl = logger.Info
	}

	logger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             opts.slowThreshold,
			LogLevel:                  lvl,
			IgnoreRecordNotFoundError: opts.ignoreNotFoundError,
			Colorful:                  true,
		},
	)

	var (
		db  *gorm.DB
		err error
	)
	switch opts.dialect {
	case MySQL:
		db, err = gorm.Open(mysql.Open(opts.MySQLDSN()), &gorm.Config{
			Logger: logger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: opts.sigularTable,
			},
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			AllowGlobalUpdate:      true,
		})
	default:
		return nil, errors.New("not support dialect")
	}

	//retry
	if err != nil {
		if db == nil {
			return nil, err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}

		var retryErr error
		var maxAttempts = 10
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			retryErr = sqlDB.Ping()
			if retryErr == nil {
				break
			}
			logrus.Warnf("ping failed,err=%v, left retry times=%d", retryErr, maxAttempts-attempts)
			time.Sleep(time.Duration(attempts) * time.Second)
		}
		if retryErr != nil {
			sqlDB.Close()
			return nil, retryErr
		}
	}

	//setMax limit
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(opts.maxOpenConn)
	sqlDB.SetMaxIdleConns(opts.maxIdleConns)
	sqlDB.SetConnMaxLifetime(opts.maxLifetime)
	err = db.Use(gormopentracing.New())
	if err != nil {
		logrus.Error("gormopentracing err ", err)
	}
	return db, nil
}
