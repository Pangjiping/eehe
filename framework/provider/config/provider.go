package config

import (
	"path/filepath"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type EeheConfigProvider struct{}

func (eehe *EeheConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewEeheConfig
}

func (eehe *EeheConfigProvider) Boot(c framework.Container) error {
	return nil
}

func (eehe *EeheConfigProvider) IsDefer() bool {
	return false
}

func (eehe *EeheConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()

	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

func (eehe *EeheConfigProvider) Name() string {
	return contract.ConfigKey
}
