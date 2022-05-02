package timeout

import "time"

type TimeoutOption func(t *Timeout)

func WithTimeout(d time.Duration) TimeoutOption {
	return func(t *Timeout) {
		t.d = d
	}
}
