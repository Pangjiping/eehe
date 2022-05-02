package ssh

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

// SSHProvider 提供eehe ssh的具体实现
type EeheSSHProvider struct{}

func (provider *EeheSSHProvider) Register(container framework.Container) framework.NewInstance {
	return NewEeheSSHService
}

func (provider *EeheSSHProvider) Boot(container framework.Container) error {
	return nil
}

func (provider *EeheSSHProvider) IsDefer() bool {
	return true
}

func (provider *EeheSSHProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (provider *EeheSSHProvider) Name() string {
	return contract.SSHKey
}
