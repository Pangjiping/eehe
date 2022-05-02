package orm

import (
	"context"
	"sync"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// EeheGorm eehe框架的orm实现
type EeheGormService struct {
	contrainer framework.Container // 服务容器
	dbs        map[string]*gorm.DB // key为dsn，value为gorm.DB(连接池)
	mu         *sync.RWMutex
}

// NewEeheGorm 实例化Gorm
func NewEeheGormService(params ...interface{}) (interface{}, error) {
	contrainer := params[0].(framework.Container)
	dbs := make(map[string]*gorm.DB)
	mu := &sync.RWMutex{}

	return &EeheGormService{
		contrainer: contrainer,
		dbs:        dbs,
		mu:         mu,
	}, nil
}

// GetDB 获取DB实例
func (svc *EeheGormService) GetDB(option ...contract.DBOption) (*gorm.DB, error) {
	logger := svc.contrainer.MustMake(contract.LogKey).(contract.Log)
	config := GetBaseConfig(svc.contrainer)
	logService := svc.contrainer.MustMake(contract.LogKey).(contract.Log)

	ormLogger := NewOrmLogger(logService)
	config.Config = &gorm.Config{
		Logger: ormLogger,
	}

	// option对opt进行修改
	for _, opt := range option {
		if err := opt(svc.contrainer, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn，就生成dsn
	if config.Dsn == "" {
		dsn, err := config.FormatDsn()
		if err != nil {
			return nil, err
		}
		config.Dsn = dsn
	}

	// 判断是否已经实例化了gorm.DB
	svc.mu.RLock()
	if db, ok := svc.dbs[config.Dsn]; ok {
		svc.mu.RUnlock()
		return db, nil
	}
	svc.mu.RUnlock()

	// 如果没有实例化gorm.DB，就进行实例化操作
	svc.mu.Lock()
	defer svc.mu.Unlock()

	// 实例化gorm.DB
	var db *gorm.DB
	var err error
	switch config.Driver {
	case "mysql":
		db, err = gorm.Open(mysql.Open(config.Dsn), config)
	case "postgres":
		db, err = gorm.Open(postgres.Open(config.Dsn), config)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.Dsn), config)
	case "sqlserver":
		db, err = gorm.Open(sqlserver.Open(config.Dsn), config)
	case "clickhouse":
		db, err = gorm.Open(clickhouse.Open(config.Dsn), config)
	}

	if err != nil {
		return nil, err
	}

	// 设置对应的连接池设置
	sqlDB, err := db.DB()
	if err != nil {
		return db, err
	}

	if config.ConnMaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.ConnMaxIdle)
	}

	if config.ConnMaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.ConnMaxOpen)
	}

	if config.ConnMaxLifetime != "" {
		liftTime, err := time.ParseDuration(config.ConnMaxLifetime)
		if err != nil {
			logger.Error(context.Background(), "conn max lift time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxLifetime(liftTime)
		}
	}

	if config.ConnMaxIdletime != "" {
		idleTime, err := time.ParseDuration(config.ConnMaxIdletime)
		if err != nil {
			logger.Error(context.Background(), "conn max idle time error", map[string]interface{}{
				"err": err,
			})
		} else {
			sqlDB.SetConnMaxIdleTime(idleTime)
		}
	}

	// 挂载到map中，结束配置
	if err == nil {
		svc.dbs[config.Dsn] = db
	}
	return db, err

}
