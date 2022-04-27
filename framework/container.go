package framework

import (
	"errors"
	"fmt"
	"sync"
)

type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换，返回error
	Bind(provider ServiceProvider) error

	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)

	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者
	Mustmake(key string) interface{}

	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

type EeheContainer struct {
	Container
	// providers 存储注册的服务提供者，key为字符串凭证
	providers map[string]ServiceProvider

	// instances 存储具体的实例，key为字符串凭证
	instances map[string]interface{}

	// mu 用于对容器的变更操作加锁
	mu sync.RWMutex
}

func NewEeheContainer() *EeheContainer {
	return &EeheContainer{
		providers: make(map[string]ServiceProvider),
		instances: make(map[string]interface{}),
		mu:        sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (eehe *EeheContainer) PrintProviders() []string {
	res := []string{}
	for _, provider := range eehe.providers {
		name := provider.Name()
		line := fmt.Sprint(name)
		res = append(res, line)
	}
	return res
}

func (eehe *EeheContainer) Bind(provider ServiceProvider) error {
	eehe.mu.Lock()
	key := provider.Name()
	eehe.providers[key] = provider
	eehe.mu.Unlock()

	// if provider is not defer
	if !provider.IsDefer() {
		if err := provider.Boot(eehe); err != nil {
			return err
		}

		params := provider.Params(eehe)
		method := provider.Register(eehe)
		instance, err := method(params...)
		if err != nil {
			fmt.Println("bind service provider ", key, " error: ", err)
			return errors.New(err.Error())
		}
		eehe.instances[key] = instance
	}
	return nil
}

func (eehe *EeheContainer) IsBind(key string) bool {
	return eehe.findServiceProvider(key) != nil
}

func (eehe *EeheContainer) findServiceProvider(key string) ServiceProvider {
	eehe.mu.Lock()
	defer eehe.mu.Unlock()

	if sp, ok := eehe.providers[key]; ok {
		return sp
	}
	return nil
}

func (eehe *EeheContainer) Make(key string) (interface{}, error) {
	return eehe.make(key, nil, false)
}

func (eehe *EeheContainer) MustMake(key string) interface{} {
	service, err := eehe.make(key, nil, false)
	if err != nil {
		panic("container not contain key " + key)
	}
	return service
}

func (eehe *EeheContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return eehe.make(key, params, true)
}

func (eehe *EeheContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(eehe); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(eehe)
	}
	method := sp.Register(eehe)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// 真正的实例化一个服务
func (eehe *EeheContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	eehe.mu.RLock()
	defer eehe.mu.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := eehe.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return eehe.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	if ins, ok := eehe.instances[key]; ok {
		return ins, nil
	}

	// 容器中还未实例化，则进行一次实例化
	inst, err := eehe.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	eehe.instances[key] = inst
	return inst, nil
}

// NameList 列出容器中所有服务提供者的字符串凭证
func (eehe *EeheContainer) NameList() []string {
	ret := []string{}
	for _, provider := range eehe.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}
