package cors

// CORSOption ...
type CORSOption func(c *CORS)

// WithAllowOrigins set allow options.
func WithAllowOrigins(allowOrigins []string) CORSOption {
	return func(c *CORS) {
		c.allowOrigins = allowOrigins
	}
}

// WithAllowMethods set allow methods.
func WithAllowMethods(allowMethods []string) CORSOption {
	return func(c *CORS) {
		c.allowMethods = allowMethods
	}
}

// WithExposeHeaders set expose headers.
func WithExposeHeaders(exposeHeaders []string) CORSOption {
	return func(c *CORS) {
		c.exposeHeaders = exposeHeaders
	}
}

// WithAllowCredentials set allow credentials.
func WithAllowCredentials(allowCredentials bool) CORSOption {
	return func(c *CORS) {
		c.allowCredentials = allowCredentials
	}
}

// WithMaxAge set max age.
func WithMaxAge(maxAge string) CORSOption {
	return func(c *CORS) {
		c.maxAge = maxAge
	}
}
