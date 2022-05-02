package timeout

import "time"

var defaultTimeout = Timeout{
	d: 10 * time.Second,
}
