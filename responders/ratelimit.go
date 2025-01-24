package responders

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"net/http"
)

type RateLimitResponder struct {
}

func (r *RateLimitResponder) ServeHTTP(w http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	req.Header.Set("X-RateLimit-Apply", "true")

	// Continue with the handler chain
	return next.ServeHTTP(w, req)
}
