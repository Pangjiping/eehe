package app

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/util"
	"github.com/google/uuid"
)

// EeheApp 代表eehe框架的app实现
type EeheApp struct {
	container  framework.Container // 服务容器
	baseFolder string              // 基础路径
	appID      string              // 表示当前这个app唯一的ID，可用于分布式锁
	configMap  map[string]string   // 配置加载
}

// AppID 获取当前app实例的uuid
func (app EeheApp) AppID() string {
	return app.appID
}

// Version 获取当前app的版本号
func (app EeheApp) Version() string {
	return "0.0.1"
}

// BaseFolder 基础目录，代表开发目录或生产运行时目录
func (app EeheApp) BaseFolder() string {
	if app.baseFolder != "" {
		return app.baseFolder
	}
	return util.GetExecDir()
}

// ConfigFolder 配置文件地址
func (app EeheApp) ConfigFolder() string {
	if val, ok := app.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "config")
}

// LogFolder 日志文件地址
func (app EeheApp) LogFolder() string {
	if val, ok := app.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "log")
}

// HttpFolder http下载数据存放地址
func (app EeheApp) HttpFolder() string {
	if val, ok := app.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "http")
}

// ConsoleFolder ...
func (app EeheApp) ConsoleFolder() string {
	if val, ok := app.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "console")
}

// StorageFolder ...
func (app EeheApp) StorageFolder() string {
	if val, ok := app.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (app EeheApp) ProviderFolder() string {
	if val, ok := app.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app", "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (app EeheApp) MiddlewareFolder() string {
	if val, ok := app.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(app.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (app EeheApp) CommandFolder() string {
	if val, ok := app.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(app.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (app EeheApp) RuntimeFolder() string {
	if val, ok := app.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(app.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (app EeheApp) TestFolder() string {
	if val, ok := app.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "test")
}

// DeployFolder 定义测试需要的信息
func (app EeheApp) DeployFolder() string {
	if val, ok := app.configMap["deploy_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "deploy")
}

// NewEeheApp 初始化EeheApp
func NewEeheApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}

	// 有两个参数，一个是容器，一个是baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
		flag.Parse()
	}
	appId := uuid.New().String()
	configMap := map[string]string{}

	return &EeheApp{
		baseFolder: baseFolder,
		container:  container,
		appID:      appId,
		configMap:  configMap}, nil
}

// LoadAppConfig 加载配置map
func (app *EeheApp) LoadAppConfig(kv map[string]string) {
	for key, val := range kv {
		app.configMap[key] = val
	}
}

// AppFolder 代表app目录
func (app *EeheApp) AppFolder() string {
	if val, ok := app.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(app.BaseFolder(), "app")
}
