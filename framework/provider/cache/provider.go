package cache

import (
	"strings"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	services "github.com/Pangjiping/eehe/framework/provider/cache/service"
)

// EeheCacheProvider EeheCache服务提供者
type EeheCacheProvider struct {
	framework.ServiceProvider
	Driver string // Driver
}

// Register 注册cache服务实例
func (ecp *EeheCacheProvider) Register(c framework.Container) framework.NewInstance {
	if ecp.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewMemoryCache
		}

		cs := tcs.(contract.Config)
		ecp.Driver = strings.ToLower(cs.GetString("cache.driver"))
	}

	switch ecp.Driver {
	case "redis":
		return services.NewRedisCache
	case "memory":
		return services.NewMemoryCache
	default:
		return services.NewMemoryCache
	}
}

// Boot 启动注入
func (ecp *EeheCacheProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer 是否延迟加载
func (ecp *EeheCacheProvider) IsDefer() bool {
	return true
}

// Params 定义要传给实例化方法的参数
func (ecp *EeheCacheProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}

// Name 定义对应的服务字符串凭证
func (ecp *EeheCacheProvider) Name() string {
	return contract.CacheKey
}
