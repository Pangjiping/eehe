package cors

import (
	"net/http"
	"strings"

	"github.com/Pangjiping/eehe/framework/gin"
)

// Func define CORS middleware.
func (m *CORS) Func() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add(varyHeader, originHeader)
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Add(varyHeader, requestMethod)
			c.Writer.Header().Add(varyHeader, requestHeaders)
		}

		origin := c.Request.Header.Get(originHeader)
		if b, allowOrigins := isOriginAllowed(m.allowOrigins, origin); b {
			m.setCORSHeader(c, allowOrigins)
		}

		if c.Request.Method == http.MethodOptions {
			c.JSON(http.StatusNoContent, nil)
			return
		}

		c.Next()
	}
}

// CORS cors
type CORS struct {
	allowOrigins     []string
	allowMethods     []string
	allowHeaders     []string
	exposeHeaders    []string
	allowCredentials bool
	maxAge           string
}

func NewCORS(opts ...CORSOption) *CORS {
	cors := defaultCORS
	for _, opt := range opts {
		opt(&cors)
	}
	return &cors
}

func (m *CORS) setCORSHeader(c *gin.Context, origins []string) {
	c.Writer.Header().Set(allowOrigin, strings.Join(origins, ","))
	c.Writer.Header().Set(allowMethods, strings.Join(m.allowMethods, ","))
	c.Writer.Header().Set(allowHeaders, strings.Join(m.allowHeaders, ","))
	c.Writer.Header().Set(exposeHeaders, strings.Join(m.exposeHeaders, ","))

	if m.allowCredentials {
		c.Writer.Header().Set(allowCredentials, "true")
	}

	if m.maxAge != "0" {
		c.Writer.Header().Set(maxAgeHeader, m.maxAge)
	}
}

func isOriginAllowed(allows []string, origin string) (bool, []string) {
	for _, o := range allows {
		if o == "*" {
			return true, []string{"*"}
		}

		if o == origin {
			return true, []string{origin}
		}
	}

	return false, []string{}
}
