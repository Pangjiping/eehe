package ratelimit

var defaultRateLimit = RateLimit{
	cap:             1000,
	rate:            500,
	waitMaxDuration: 0,
}
