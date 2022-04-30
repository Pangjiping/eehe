package kernel

import (
	"net/http"

	"github.com/Pangjiping/eehe/framework/gin"
)

// EeheKernelService provides engine service.
type EeheKernelService struct {
	engine *gin.Engine
}

// NewEeheKernelService init web service engine(gin).
func NewEeheKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &EeheKernelService{
		engine: httpEngine}, nil
}

// HttpEngine gets web engine.
func (service *EeheKernelService) HttpEngine() http.Handler {
	return service.engine
}
