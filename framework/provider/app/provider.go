package app

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

// EeheAppProvider 提供app的具体实现方法
type EeheAppProvider struct {
	BaseFolder string
}

// Register 注册HadeApp方法
func (e *EeheAppProvider) Register(container framework.Container) framework.NewInstance {
	return NewEeheApp
}

// Boot 启动调用
func (e *EeheAppProvider) Boot(container framework.Container) error {
	return nil
}

// IsDefer 是否延迟初始化
func (e *EeheAppProvider) IsDefer() bool {
	return false
}

// Params 获取初始化参数
func (e *EeheAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, e.BaseFolder}
}

// Name 获取字符串凭证
func (e *EeheAppProvider) Name() string {
	return contract.AppKey
}
