package distributed

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type LocalDistributedProvider struct{}

func (ldp *LocalDistributedProvider) Register(container framework.Container) framework.NewInstance {
	return NewLocalDistributedService
}

func (ldp *LocalDistributedProvider) Boot(container framework.Container) error {
	return nil
}

func (ldp *LocalDistributedProvider) IsDefer() bool {
	return false
}

func (ldp *LocalDistributedProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container}
}

func (ldp *LocalDistributedProvider) Name() string {
	return contract.DistributedKey
}
