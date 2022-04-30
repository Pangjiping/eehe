package kernel

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/gin"
)

type EeheKernelProvider struct {
	HttpEngine *gin.Engine
}

func (provider *EeheKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewEeheKernelService
}

func (provider *EeheKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}
	provider.HttpEngine.SetContainer(c)
	return nil
}

func (provider *EeheKernelProvider) IsDefer() bool {
	return false
}

func (provider *EeheKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

func (provider *EeheKernelProvider) Name() string {
	return contract.KernelKey
}
