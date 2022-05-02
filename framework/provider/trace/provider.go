package trace

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type EeheTraceProvider struct {
	c framework.Container
}

func (provider *EeheTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewEeheTraceService
}

func (provider *EeheTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

func (provider *EeheTraceProvider) IsDefer() bool {
	return false
}

func (provider *EeheTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

func (provider *EeheTraceProvider) Name() string {
	return contract.TraceKey
}
