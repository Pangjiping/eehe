package services

import (
	"context"
	"io"
	pkgLog "log"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/provider/log/formatter"
)

type EeheLog struct {
	level      contract.LogLevel   // log level
	formatter  contract.Formatter  // formatter
	ctxFielder contract.CtxFielder // context
	output     io.Writer           // output
	c          framework.Container // app container
}

func (log *EeheLog) IsLevelEnable(level contract.LogLevel) bool {
	return level <= log.level
}

func (log *EeheLog) logf(level contract.LogLevel, ctx context.Context, msg string, fields map[string]interface{}) error {
	if !log.IsLevelEnable(level) {
		return nil
	}

	fs := fields
	if log.ctxFielder != nil {
		t := log.ctxFielder(ctx)
		if t != nil {
			for k, v := range t {
				fs[k] = v
			}
		}
	}

	if log.c != nil && log.c.IsBind(contract.TraceKey) {
		tracer := log.c.MustMake(contract.TraceKey).(contract.Trace)
		tc := tracer.GetTrace(ctx)
		if tc != nil {
			maps := tracer.ToMap(tc)
			for k, v := range maps {
				fs[k] = v
			}
		}
	}

	if log.formatter == nil {
		log.formatter = formatter.TextFormatter
	}

	ct, err := log.formatter(level, time.Now(), msg, fs)
	if err != nil {
		return err
	}

	if level == contract.PanicLevel {
		pkgLog.Panicln(string(ct))
		return nil
	}

	log.output.Write(ct)
	log.output.Write([]byte("\r\n"))
	return nil
}

func (log *EeheLog) SetOutput(output io.Writer) {
	log.output = output
}

func (log *EeheLog) Panic(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.PanicLevel, ctx, msg, fields)
}

func (log *EeheLog) Fatal(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.FatalLevel, ctx, msg, fields)
}

func (log *EeheLog) Error(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.ErrorLevel, ctx, msg, fields)
}

func (log *EeheLog) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.WarnLevel, ctx, msg, fields)
}

func (log *EeheLog) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.InfoLevel, ctx, msg, fields)
}

func (log *EeheLog) Debug(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.DebugLevel, ctx, msg, fields)
}

func (log *EeheLog) Trace(ctx context.Context, msg string, fields map[string]interface{}) {
	log.logf(contract.TraceLevel, ctx, msg, fields)
}

func (log *EeheLog) SetLevel(level contract.LogLevel) {
	log.level = level
}

func (log *EeheLog) SetCtxFielder(handler contract.CtxFielder) {
	log.ctxFielder = handler
}

func (log *EeheLog) SetFormatter(formatter contract.Formatter) {
	log.formatter = formatter
}
