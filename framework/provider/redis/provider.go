package redis

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

// RedisProvider 提供redis的具体实现方法
type EeheRedisProvider struct{}

func (provider *EeheRedisProvider) Register(container framework.Container) framework.NewInstance {
	return NewEeheRedisService
}

func (provider *EeheRedisProvider) Boot(container framework.Container) error {
	return nil
}

func (provider *EeheRedisProvider) IsDefer() bool {
	return true
}

func (provider *EeheRedisProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (provider *EeheRedisProvider) Name() string {
	return contract.RedisKey
}
