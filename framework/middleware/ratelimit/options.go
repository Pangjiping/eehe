package ratelimit

import "time"

type RateLimitOption func(r *RateLimit)

// WithRate set rate.
func WithRate(rate float64) RateLimitOption {
	return func(r *RateLimit) {
		r.rate = rate
	}
}

// WithCap set cap.
func WithCap(cap int64) RateLimitOption {
	return func(r *RateLimit) {
		r.cap = cap
	}
}

// WithWaitMaxDuration set wait max duration.
func WithWaitMaxDuration(max time.Duration) RateLimitOption {
	return func(r *RateLimit) {
		r.waitMaxDuration = max
	}
}
