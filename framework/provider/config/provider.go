package config

import (
	"path/filepath"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type EeheConfigProvider struct{}

// Register registe a new function for make a service instance.
func (eehe *EeheConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewEeheConfig
}

// Boot will called when the service instantiate.
func (eehe *EeheConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register.
func (eehe *EeheConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance.
func (eehe *EeheConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()

	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

// Name defines the name for this service.
func (eehe *EeheConfigProvider) Name() string {
	return contract.ConfigKey
}
