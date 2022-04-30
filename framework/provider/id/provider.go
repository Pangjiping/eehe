package id

import (
	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
)

type EeheIDProvider struct{}

// Register registes a new function for make a service instance.
func (provider *EeheIDProvider) Register(c framework.Container) framework.NewInstance {
	return NewEeheIDService
}

// Boot will be called when the service instantiate.
func (provider *EeheIDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer defines whether the service instantiate when first make or register.
func (provider *EeheIDProvider) IsDefer() bool {
	return false
}

// Params defines the necessary params for NewInstance.
func (provider *EeheIDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

// Name defines the name of this service.
func (provider *EeheIDProvider) Name() string {
	return contract.IDKey
}
