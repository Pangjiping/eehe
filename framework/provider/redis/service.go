package redis

import (
	"sync"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/go-redis/redis/v8"
)

// EeheRedis 代表eehe框架的redis实现
type EeheRedisService struct {
	container framework.Container      // 服务容器
	clients   map[string]*redis.Client // key为uniqueKey，value为redisClient(连接池)
	mu        *sync.RWMutex
}

// NewEeheRedis 实例化redis client
func NewEeheRedisService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*redis.Client)
	mu := &sync.RWMutex{}

	return &EeheRedisService{
		container: container,
		clients:   clients,
		mu:        mu,
	}, nil
}

func (svc *EeheRedisService) GetClient(options ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetBaseConfig(svc.container)

	// option对opt进行修改
	for _, opt := range options {
		if err := opt(svc.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn，就生成dsn
	key := config.UniqKey()

	// 判断是否已经实例化了redis.Client
	svc.mu.RLock()
	if db, ok := svc.clients[key]; ok {
		svc.mu.RUnlock()
		return db, nil
	}
	svc.mu.RLock()

	// 如果没有实现redis，就要实例化
	svc.mu.Lock()
	defer svc.mu.Unlock()

	client := redis.NewClient(config.Options)
	svc.clients[key] = client

	return client, nil
}
