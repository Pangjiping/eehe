package orm

import (
	"context"
	"time"

	"github.com/Pangjiping/eehe/framework/contract"
	"gorm.io/gorm/logger"
)

// OrmLogger orm日志的实现类，实现了gorm.Logger.Interface
type OrmLogger struct {
	logger contract.Log // 有一个logger对象存放eehe的log服务
}

func NewOrmLogger(logger contract.Log) *OrmLogger {
	return &OrmLogger{
		logger: logger,
	}
}

func (orm *OrmLogger) LogMode(level logger.LogLevel) logger.Interface {
	return orm
}

func (orm *OrmLogger) Info(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	orm.logger.Info(ctx, s, fields)
}

func (orm *OrmLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	orm.logger.Warn(ctx, s, fields)
}

func (orm *OrmLogger) Error(ctx context.Context, s string, i ...interface{}) {
	fields := map[string]interface{}{
		"fields": i,
	}
	orm.logger.Error(ctx, s, fields)
}

func (orm *OrmLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	elapsed := time.Since(begin)
	fields := map[string]interface{}{
		"begin": begin,
		"error": err,
		"sql":   sql,
		"rows":  rows,
		"time":  elapsed,
	}

	s := "orm trace sql"
	orm.logger.Trace(ctx, s, fields)
}
