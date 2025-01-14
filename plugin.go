package caddydefender

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/jasonlovesdoggo/caddy-defender/responders"
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
	Responder Responder `json:"responder,omitempty"`
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

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (m *DefenderMiddleware) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// Parse additional IP ranges
		for d.NextArg() {
			m.AdditionalRanges = append(m.AdditionalRanges, d.Val())
		}

		// Parse responder backend
		if d.NextArg() {
			switch d.Val() {
			case "block":
				m.Responder = responders.BlockResponder{}
			case "garbage":
				m.Responder = responders.GarbageResponder{}
			case "custom":
				if !d.NextArg() {
					return d.ArgErr()
				}
				m.Responder = responders.CustomResponder{Message: d.Val()}
			default:
				return d.Errf("unknown responder: %s", d.Val())
			}
		}
	}
	return nil
}
