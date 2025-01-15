package caddydefender

import (
	"encoding/json"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	// Register the module with Caddy
	caddy.RegisterModule(DefenderMiddleware{})
	httpcaddyfile.RegisterHandlerDirective("defender", parseCaddyfile)
}

// DefenderMiddleware implements an HTTP middleware that enforces IP-based rules.
type DefenderMiddleware struct {
	// Additional IP ranges specified by the user
	AdditionalRanges []string `json:"additional_ranges,omitempty"`
	// Responder backend to use
	// Use concrete responder type for JSON
	ResponderRaw json.RawMessage `json:"responder,omitempty"`
	// Internal field for the actual responder interface
	responder       Responder       `json:"-"`
	ResponderConfig json.RawMessage `json:"responder_config,omitempty"`

	// Logger
	log *zap.Logger
}

// Provision sets up the middleware and logger.
func (m *DefenderMiddleware) Provision(ctx caddy.Context) error {
	m.log = ctx.Logger(m)
	return nil
}

// CaddyModule returns the Caddy module information.
func (DefenderMiddleware) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.defender",
		New: func() caddy.Module { return new(DefenderMiddleware) },
	}
}

// parseCaddyfile parses the Caddyfile directive.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m DefenderMiddleware
	err := m.UnmarshalCaddyfile(h.Dispenser)
	if err != nil {
		return nil, err
	}
	return m, nil
}
