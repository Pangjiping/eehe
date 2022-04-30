package formatter

import "github.com/Pangjiping/eehe/framework/contract"

func Prefix(level contract.LogLevel) string {
	prefix := ""
	switch level {
	case contract.PanicLevel:
		prefix = "[PANIC]"
	case contract.FatalLevel:
		prefix = "[FATAL]"
	case contract.ErrorLevel:
		prefix = "[ERROR]"
	case contract.WarnLevel:
		prefix = "[WARNNING]"
	case contract.InfoLevel:
		prefix = "[INFO]"
	case contract.DebugLevel:
		prefix = "[DEBUG]"
	case contract.TraceLevel:
		prefix = "[TRACE]"
	}
	return prefix
}
