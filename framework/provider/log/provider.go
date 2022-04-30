package log

import (
	"io"
	"strings"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/provider/log/formatter"
	"github.com/Pangjiping/eehe/framework/provider/log/services"
)

type EeheLogServiceProvider struct {
	framework.ServiceProvider
	Driver     string // driver
	Level      contract.LogLevel
	Formatter  contract.Formatter
	CtxFielder contract.CtxFielder
	Output     io.Writer
}

func (provider *EeheLogServiceProvider) Register(c framework.Container) framework.NewInstance {
	if provider.Driver == "" {
		tcs, err := c.Make(contract.ConfigKey)
		if err != nil {
			return services.NewEeheConsoleLog
		}

		cs := tcs.(contract.Config)
		provider.Driver = strings.ToLower(cs.GetString("log.Driver"))
	}

	switch provider.Driver {
	case "single":
		return services.NewEeheSingleLog
	case "rotate":
		return services.NewEeheRotateLog
	case "console":
		return services.NewEeheConsoleLog
	case "custom":
		return services.NewEeheCustomLog
	default:
		return services.NewEeheConsoleLog
	}
}

func (provider *EeheLogServiceProvider) Boot(c framework.Container) error {
	return nil
}

func (provider *EeheLogServiceProvider) IsDefer() bool {
	return false
}

func (provider *EeheLogServiceProvider) Params(c framework.Container) []interface{} {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)

	if provider.Formatter == nil {
		provider.Formatter = formatter.TextFormatter
		if configService.IsExist("log.formatter") {
			v := configService.GetString("log.formatter")
			if v == "json" {
				provider.Formatter = formatter.JsonFormatter
			} else if v == "text" {
				provider.Formatter = formatter.TextFormatter
			}
		}
	}

	if provider.Level == contract.UnknownLevel {
		provider.Level = contract.InfoLevel
		if configService.IsExist("log.level") {
			provider.Level = logLevel(configService.GetString("log.level"))
		}
	}

	return []interface{}{c, provider.Level, provider.CtxFielder, provider.Formatter, provider.Output}
}

func (provider *EeheLogServiceProvider) Name() string {
	return contract.LogKey
}

func logLevel(config string) contract.LogLevel {
	switch strings.ToLower(config) {
	case "panic":
		return contract.PanicLevel
	case "info":
		return contract.InfoLevel
	case "warn":
		return contract.WarnLevel
	case "fatal":
		return contract.FatalLevel
	case "error":
		return contract.ErrorLevel
	case "debug":
		return contract.DebugLevel
	case "trace":
		return contract.TraceLevel
	default:
	}
	return contract.UnknownLevel
}
