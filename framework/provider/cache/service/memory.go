/*
	实现一个内存缓存
	包括了简单的增删改查和过期操作，支持远程调用获取分布式锁服务
	map+RWMutex实现，不适合高并发写操作
	TODO: 提供一个并发写的MemoryCache
*/

package services

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type MemoryData struct {
	val       interface{}
	createdAt time.Time
	ttl       time.Duration
}

type MemoryCache struct {
	container framework.Container
	data      map[string]*MemoryData
	mu        sync.RWMutex
}

func NewMemoryCache(opts ...interface{}) (interface{}, error) {
	c := opts[0].(framework.Container)
	obj := &MemoryCache{
		container: c,
		data:      make(map[string]*MemoryData),
		mu:        sync.RWMutex{},
	}
	return obj, nil
}

func (mc *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	var val string
	if err := mc.GetObj(ctx, key, &val); err != nil {
		return "", err
	}
	return val, nil
}

func (mc *MemoryCache) GetObj(ctx context.Context, key string, obj interface{}) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if md, ok := mc.data[key]; ok {
		if md.ttl != NoneDuration {
			if time.Now().Sub(md.createdAt) > md.ttl {
				delete(mc.data, key)
				return ErrKeyNotFound
			}
		}

		bt, _ := json.Marshal(md.val)
		err := json.Unmarshal(bt, obj)
		if err != nil {
			return err
		}
		return nil
	}

	return ErrKeyNotFound
}

func (mc *MemoryCache) GetMany(ctx context.Context, keys []string) (map[string]string, error) {
	errs := make([]string, 0, len(keys))
	rets := make(map[string]string)
	for _, key := range keys {
		val, err := mc.Get(ctx, key)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		rets[key] = val
	}
	if len(errs) == 0 {
		return rets, nil
	}
	return rets, errors.New(strings.Join(errs, "||"))
}

func (mc *MemoryCache) Set(ctx context.Context, key string, val string, timeout time.Duration) error {
	return mc.SetObj(ctx, key, val, timeout)
}

func (mc *MemoryCache) SetObj(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	md := &MemoryData{
		val:       val,
		createdAt: time.Now(),
		ttl:       timeout,
	}
	mc.data[key] = md
	return nil
}

func (mc *MemoryCache) SetMany(ctx context.Context, data map[string]string, timeout time.Duration) error {
	errs := []string{}
	for k, v := range data {
		err := mc.Set(ctx, k, v, timeout)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "||"))
	}
	return nil
}

func (mc *MemoryCache) SetForever(ctx context.Context, key string, val string) error {
	return mc.Set(ctx, key, val, NoneDuration)
}

func (mc *MemoryCache) SetForeverObj(ctx context.Context, key string, val interface{}) error {
	return mc.SetObj(ctx, key, val, NoneDuration)
}

func (mc *MemoryCache) Remember(ctx context.Context, key string, timeout time.Duration, rememberFunc contract.RememberFunc, obj interface{}) error {
	err := mc.GetObj(ctx, key, obj)
	if err == nil {
		return nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return err
	}

	// key not found
	objNew, err := rememberFunc(ctx, mc.container)
	if err != nil {
		return err
	}

	if err := mc.SetObj(ctx, key, objNew, timeout); err != nil {
		return err
	}

	if err := mc.GetObj(ctx, key, &obj); err != nil {
		return err
	}
	return nil
}

func (mc *MemoryCache) SetTTL(ctx context.Context, key string, timeout time.Duration) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	if md, ok := mc.data[key]; ok {
		md.ttl = timeout
		return nil
	}
	return ErrKeyNotFound
}

func (mc *MemoryCache) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	if md, ok := mc.data[key]; ok {
		return md.ttl, nil
	}
	return NoneDuration, ErrKeyNotFound
}

func (mc *MemoryCache) Calc(ctx context.Context, key string, step int64) (int64, error) {
	var val int64
	err := mc.GetObj(ctx, key, &val)
	val = val + step
	if err == nil {
		mc.data[key].val = val
		return val, nil
	}

	if !errors.Is(err, ErrKeyNotFound) {
		return 0, err
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()
	// key not found
	mc.data[key] = &MemoryData{
		val:       val,
		createdAt: time.Now(),
		ttl:       NoneDuration,
	}

	return val, nil
}

func (mc *MemoryCache) Increment(ctx context.Context, key string) (int64, error) {
	return mc.Calc(ctx, key, 1)
}

func (mc *MemoryCache) Decrement(ctx context.Context, key string) (int64, error) {
	return mc.Calc(ctx, key, -1)
}

func (mc *MemoryCache) Del(ctx context.Context, key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.data, key)
	return nil
}

func (mc *MemoryCache) DelMany(ctx context.Context, keys []string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	for _, key := range keys {
		delete(mc.data, key)
	}
	return nil
}
