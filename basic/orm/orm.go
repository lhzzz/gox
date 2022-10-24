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

func logrusLevelToOrmLevel(logruslevel logrus.Level) logger.LogLevel {
	return logger.Info
}

func OpenDB(opt ...Option) (*gorm.DB, error) {
	opts := defaultDbOptions
	for _, o := range opt {
		o(&opts)
	}

	logger := logger.New(
		logrus.StandardLogger(),
		logger.Config{
			SlowThreshold:             opts.slowThreshold,
			LogLevel:                  logrusLevelToOrmLevel(logrus.GetLevel()),
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
			return nil, errors.New("db connpool is not a ping interface")
		}

		var dbError error
		maxAttempts := 10
		for attempts := 1; attempts <= maxAttempts; attempts++ {
			dbError := sqlDB.Ping()
			if dbError == nil {
				break
			}
			logrus.Warnf("ping failed,err=%v, left retry times=%d", dbError, maxAttempts-attempts)
			time.Sleep(time.Duration(attempts) * time.Second)
		}
		if dbError != nil {
			sqlDB.Close()
			return nil, dbError
		}
	}

	err = db.Use(gormopentracing.New())
	if err != nil {
		logrus.Error("gormopentracing err ", err)
	}
	return db, nil
}
