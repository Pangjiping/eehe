package id

import "github.com/rs/xid"

type EeheIDService struct{}

func NewEeheIDService(params ...interface{}) (interface{}, error) {
	return &EeheIDService{}, nil
}

func (id *EeheIDService) NewID() string {
	return xid.New().String()
}
