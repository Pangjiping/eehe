package env

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type EeheEnvProvider struct {
	Folder string
}

// Register registers a new function for make a service instance.
func (provider *EeheEnvProvider) Register(container framework.Container) framework.NewInstance {
	return NewEeheEnvService
}

// Boot will be called when the service instantiate.
func (provider *EeheEnvProvider) Boot(container framework.Container) error {
	app := container.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer defines whether the service instantiate when first make or register.
func (provider *EeheEnvProvider) IsDefer() bool {
	return false
}

// Params defines the neccessary params for NewInstance.
func (provider *EeheEnvProvider) Params(container framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

// Name defines the name for this service.
func (provider *EeheEnvProvider) Name() string {
	return contract.EnvKey
}
