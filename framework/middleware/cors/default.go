package cors

import "net/http"

var defaultCORS = CORS{
	allowCredentials: false,
	allowHeaders: []string{
		"Content-Type",
		"Origin",
		"X-CSRF-Token",
		"Authorization",
		"AccessToken",
		"Token",
		"Range",
	},
	maxAge:       "86400",
	allowOrigins: []string{"*"},
	exposeHeaders: []string{
		"Content-Length",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
	},
	allowMethods: []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	},
}
