package orm

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

// GormProvider 提供orm的具体实现方法
type EeheGormProvider struct{}

// Register 注册方法
func (provider *EeheGormProvider) Register(container framework.Container) framework.NewInstance {
	return NewEeheGormService
}

// Boot 启动调用注入
func (provider *EeheGormProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (provider *EeheGormProvider) IsDefer() bool {
	return true
}

// Params 获取初始化参数
func (provider *EeheGormProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

// Name 获取字符串凭证
func (provider *EeheGormProvider) Name() string {
	return contract.ORMKey
}
