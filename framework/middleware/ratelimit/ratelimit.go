package ratelimit

import (
	"net/http"
	"time"

	"github.com/Pangjiping/eehe/framework/gin"
	"github.com/juju/ratelimit"
)

// Func define ratelimit middleware.
func (r *RateLimit) Func() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := r.bucket.TakeMaxDuration(1, r.waitMaxDuration); !ok {
			ctx.JSON(http.StatusTooManyRequests, nil)
			return
		}

		ctx.Next()
	}
}

type RateLimit struct {
	cap             int64
	rate            float64
	waitMaxDuration time.Duration
	bucket          *ratelimit.Bucket
}

func NewRateLimit(opts ...RateLimitOption) *RateLimit {
	r := defaultRateLimit
	for _, opt := range opts {
		opt(&r)
	}

	r.bucket = ratelimit.NewBucketWithRate(r.rate, r.cap)

	return &r
}
