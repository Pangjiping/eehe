package trace

import (
	"context"
	"net/http"
	"time"

	"github.com/Pangjiping/eehe/framework"
	"github.com/Pangjiping/eehe/framework/contract"
	"github.com/Pangjiping/eehe/framework/gin"
)

type TraceKey string

var ContextKey = TraceKey("trace-key")

type EeheTraceService struct {
	idService        contract.IDService
	traceIDGenerator contract.IDService
	spanIDGenerator  contract.IDService
}

// NewEeheTraceService 获取trace实例
func NewEeheTraceService(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	idService := c.MustMake(contract.IDKey).(contract.IDService)
	return &EeheTraceService{
		idService: idService,
	}, nil
}

// WithTrace register new trace to context.
func (svc *EeheTraceService) WithTrace(ctx context.Context, trace *contract.TraceContext) context.Context {
	if ginC, ok := ctx.(*gin.Context); ok {
		ginC.Set(string(ContextKey), trace)
		return ginC
	} else {
		newC := context.WithValue(ctx, ContextKey, trace)
		return newC
	}
}

// GetTrace form trace context.
func (svc *EeheTraceService) GetTrace(ctx context.Context) *contract.TraceContext {
	if ginC, ok := ctx.(*gin.Context); ok {
		if val, ok2 := ginC.Get(string(ContextKey)); ok2 {
			return val.(*contract.TraceContext)
		}
	}

	if tc, ok := ctx.Value(ContextKey).(*contract.TraceContext); ok {
		return tc
	}
	return nil
}

// NewTrace generate a new trace.
func (svc *EeheTraceService) NewTrace() *contract.TraceContext {
	var traceID, spanID string
	if svc.traceIDGenerator != nil {
		traceID = svc.traceIDGenerator.NewID()
	} else {
		traceID = svc.idService.NewID()
	}

	if svc.spanIDGenerator != nil {
		spanID = svc.spanIDGenerator.NewID()
	} else {
		spanID = svc.idService.NewID()
	}

	tc := &contract.TraceContext{
		TraceID:    traceID,
		ParentID:   "",
		SpanID:     spanID,
		CspanID:    "",
		Annotation: map[string]string{},
	}

	return tc
}

// StartSpan instance a sub trace with new span id.
func (svc *EeheTraceService) StartSpan(tc *contract.TraceContext) *contract.TraceContext {
	var childSpanID string
	if svc.spanIDGenerator != nil {
		childSpanID = svc.spanIDGenerator.NewID()
	} else {
		childSpanID = svc.idService.NewID()
	}

	childSpan := &contract.TraceContext{
		TraceID:  tc.TraceID,
		ParentID: "",
		SpanID:   tc.SpanID,
		CspanID:  childSpanID,
		Annotation: map[string]string{
			contract.TraceKeyTime: time.Now().String(),
		},
	}

	return childSpan
}

// ExtractHTTP get trace by http.
func (svc *EeheTraceService) ExtractHTTP(req *http.Request) *contract.TraceContext {
	tc := &contract.TraceContext{}
	tc.TraceID = req.Header.Get(contract.TraceKeyTraceID)
	tc.ParentID = req.Header.Get(contract.TraceKeySpanID)
	tc.SpanID = req.Header.Get(contract.TraceKeyCspanID)
	tc.CspanID = ""

	if tc.TraceID == "" {
		tc.TraceID = svc.idService.NewID()
	}

	if tc.SpanID == "" {
		tc.SpanID = svc.idService.NewID()
	}

	return tc
}

// InjectHTTP inject trace to http.
func (svc *EeheTraceService) InjectHTTP(req *http.Request, tc *contract.TraceContext) *http.Request {
	req.Header.Add(contract.TraceKeyTraceID, tc.TraceID)
	req.Header.Add(contract.TraceKeySpanID, tc.SpanID)
	req.Header.Add(contract.TraceKeyCspanID, tc.CspanID)
	req.Header.Add(contract.TraceKeyParentID, tc.ParentID)
	return req
}

func (svc *EeheTraceService) ToMap(tc *contract.TraceContext) map[string]string {
	m := map[string]string{}
	if tc == nil {
		return m
	}

	m[contract.TraceKeyTraceID] = tc.TraceID
	m[contract.TraceKeySpanID] = tc.SpanID
	m[contract.TraceKeyCspanID] = tc.CspanID
	m[contract.TraceKeyParentID] = tc.ParentID

	if tc.Annotation != nil {
		for k, v := range tc.Annotation {
			m[k] = v
		}
	}

	return m
}
