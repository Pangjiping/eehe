package trace

import (
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/gin"
)

// Func define trace middleware.
func (t *Trace) Func() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tracer := ctx.MustMake(contract.TraceKey).(contract.Trace)
		traceCtx := tracer.ExtractHTTP(ctx.Request)
		tracer.WithTrace(ctx, traceCtx)

		ctx.Next()
	}
}

type Trace struct{}

func NewTrace() *Trace {
	return &Trace{}
}
