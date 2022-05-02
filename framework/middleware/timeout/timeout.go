package timeout

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Pangjiping/eehe/framework/gin"
)

func (t *Timeout) Func() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 执行业务逻辑前预操作；初始化超时context
		durationCtx, cancel := context.WithTimeout(ctx.BaseContext(), t.d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()

			ctx.Next()

			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			ctx.ISetStatus(500).IJson("time out")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			ctx.ISetStatus(500).IJson("time out")
		}

	}
}

type Timeout struct {
	d time.Duration
}

func NewTimeout(opts ...TimeoutOption) *Timeout {
	t := defaultTimeout
	for _, opt := range opts {
		opt(&t)
	}
	return &t
}
