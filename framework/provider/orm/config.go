package orm

import (
	"context"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

// GetBaseConfig loads .yaml file struct.
func GetBaseConfig(c framework.Container) *contract.DBConfig {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.DBConfig{}

	// load .yaml file
	err := configService.Load("database", config)
	if err != nil {
		logService.Error(context.Background(), "parse database config error", nil)
		return nil
	}
	return config
}

// WithConfigPath loads setting file.
func WithConfigPath(configPath string) contract.DBOption {
	return func(c framework.Container, config *contract.DBConfig) error {
		configService := c.MustMake(contract.ConfigKey).(contract.Config)

		// load configPath
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}
}

// WithGormConfig 表示自行配置Gorm配置信息
func WithGormConfig(f func(options *contract.DBConfig)) contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		f(config)
		return nil
	}
}

// WithDryRun 设置空跑模式
func WithDryRun() contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

// WithFullSaveAssociations 设置保存时候关联
func WithFullSaveAssociations() contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		config.FullSaveAssociations = true
		return nil
	}
}
