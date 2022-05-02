package ssh

import (
	"context"
	"sync"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"golang.org/x/crypto/ssh"
)

// EeheSSH 代表eehe框架的ssh实现
type EeheSSHService struct {
	container framework.Container    // 服务容器
	clients   map[string]*ssh.Client // key为uniqueKey，value为ssh.Client(连接池)
	mu        *sync.RWMutex
}

// NewEeheSSH 代表实例化Client
func NewEeheSSHService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	clients := make(map[string]*ssh.Client)
	mu := &sync.RWMutex{}
	return &EeheSSHService{
		container: container,
		clients:   clients,
		mu:        mu,
	}, nil
}

func (svc *EeheSSHService) GetClient(option ...contract.SSHOption) (*ssh.Client, error) {
	logService := svc.container.MustMake(contract.LogKey).(contract.Log)
	// 读取默认值
	config := GetBaseConfig(svc.container)

	// option对opt进行修改
	for _, opt := range option {
		if err := opt(svc.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn，就生成dsn
	key := config.UniqKey()

	// 判断是否已经实例化了
	svc.mu.RLock()
	if ssh, ok := svc.clients[key]; ok {
		svc.mu.RUnlock()
		return ssh, nil
	}
	svc.mu.RUnlock()

	// 没有实例化，就要进行实例化
	svc.mu.Lock()
	defer svc.mu.Unlock()

	// 实例化
	addr := config.Host + ":" + config.Port
	client, err := ssh.Dial(config.NetWork, addr, config.ClientConfig)
	if err != nil {
		logService.Error(context.Background(), "ssh dial error", map[string]interface{}{
			"err":  err,
			"addr": addr,
		})
	}

	// 挂载到map中，结束配置
	svc.clients[key] = client
	return client, nil
}
