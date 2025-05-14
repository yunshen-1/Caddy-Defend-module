package responders

import (
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// Responder defines the interface for handling responses.
type Responder interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error
}
