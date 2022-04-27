/*
	提供缓存中间件redis客户端
*/

package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/provider/redis"
	redisv8 "github.com/go-redis/redis/v8"
)

type RedisCache struct {
	container framework.Container
	client    *redisv8.Client
	mu        sync.RWMutex
}

func NewRedisCache(opts ...interface{}) (interface{}, error) {
	container := opts[0].(framework.Container)
	if !container.IsBind(contract.RedisKey) {
		err := container.Bind(&redis.RedisProvider{})
		if err != nil {
			return nil, err
		}
	}

	redisService := container.MustMake(contract.RedisKey).(contract.RedisService)
	client, err := redisService.GetClient(redis.WithConfigPath("cache"))
	if err != nil {
		return nil, err
	}

	obj := &RedisCache{
		container: container,
		client:    client,
		mu:        sync.RWMutex{},
	}
	return obj, nil
}

func (rc *RedisCache) Get(ctx context.Context, key string) (string, error) {
	cmd := rc.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisv8.Nil) {
		return "", ErrKeyNotFound
	}
	return cmd.Result()
}

func (rc *RedisCache) GetObj(ctx context.Context, key string, model interface{}) error {
	cmd := rc.client.Get(ctx, key)
	if errors.Is(cmd.Err(), redisv8.Nil) {
		return ErrKeyNotFound
	}

	err := cmd.Scan(model)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RedisCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	pipeline := rc.client.Pipeline()
	vals := make(map[string]string)
	cmds := make([]*redisv8.StringCmd, 0, len(keys))

	for _, key := range keys {
		cmds = append(cmds, pipeline.Get(ctx, key))
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	errs := make([]string, 0, len(keys))
	for _, cmd := range cmds {
		val, err := cmd.Result()
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		key := cmd.Args()[1].(string)
		vals[key] = val
	}

	return vals, nil
}

func (rc *RedisCache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return rc.client.Set(ctx, key, val, timeout).Err()
}

func (rc *RedisCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return rc.client.Set(ctx, key, val, timeout).Err()
}

func (rc *RedisCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	pipline := rc.client.Pipeline()
	cmds := make([]*redisv8.StatusCmd, 0, len(data))
	for k, v := range data {
		cmds = append(cmds, pipline.Set(ctx, k, v, timeout))
	}
	_, err := pipline.Exec(ctx)
	return err
}

func (rc *RedisCache) SetForever(ctx context.Context, key string, val string) error {
	return rc.client.Set(ctx, key, val, NoneDuration).Err()
	return nil
}

func (rc *RedisCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return rc.client.Set(ctx, key, val, NoneDuration).Err()
}

func (rc *RedisCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	return rc.client.Expire(ctx, key, timeout).Err()
}

func (rc *RedisCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	return rc.client.TTL(ctx, key).Result()
}

func (rc *RedisCache) Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc RememberFunc, model interface{}) error {
	err := rc.GetObj(ctx, key, model)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	// key not found
	objNew, err := rememberFunc(ctx, rc.container)
	if err != nil {
		return err
	}

	if err := rc.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}
	if err := rc.GetObj(ctx, key, model); err != nil {
		return err
	}
	return nil
}

func (rc *RedisCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	return rc.client.IncrBy(ctx, key, step).Result()
}

func (rc *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	return rc.client.IncrBy(ctx, key, 1).Result()
}

func (rc *RedisCache) Decrement(ctx context.Context, key string) (int64, error) {
	return rc.client.IncrBy(ctx, key, -1).Result()
}

func (rc *RedisCache) Del(ctx context.Context, key string) error {
	return rc.client.Del(ctx, key).Err()
}

func (rc *RedisCache) DelMany(ctx context.Context, keys []string) error {
	pipline := rc.client.Pipeline()
	cmds := make([]*redisv8.IntCmd, 0, len(keys))
	for _, key := range keys {
		cmds = append(cmds, pipline.Del(ctx, key))
	}
	_, err := pipline.Exec(ctx)
	return err
}
