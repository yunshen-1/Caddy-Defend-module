package caddydefender

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"net/http"
)

// Responder defines the interface for handling responses.
type Responder interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error
}
