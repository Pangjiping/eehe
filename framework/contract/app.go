package contract

// AppKey defines string token.
const AppKey = "eehe:app"

type App interface {
	// AppID indicates ID of each app instance,
	// which will be used in scenarios such as distributed locks.
	AppID() string

	// Version returns version information of app.
	Version() string

	// BaseFolder defines address of current project.
	BaseFolder() string

	// ConfigFolder defines address of configure file.
	ConfigFolder() string

	// LogFolder defines address of log file.
	LogFolder() string

	// ProviderFolder defines address of service provider.
	ProviderFolder() string

	// MiddlewareFolder defines custom middleware.
	MiddlewareFolder() string

	// CommandFolder defines the commands of app.
	CommandFolder() string

	// RuntimeFolder defines the runtime information.
	RuntimeFolder() string

	// TestFolder defines the address of test data.
	TestFolder() string

	// DeployFolder defines deploy information.
	DeployFolder() string

	// AppFolder defines address of go code.
	AppFolder() string

	// LoadAppConfig loads new AppConfig.
	// Key is above methods, ConfigFolder -> config_folder.
	LoadAppConfig(kv map[string]string)
}
