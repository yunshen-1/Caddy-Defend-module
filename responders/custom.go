package responders

import (
	"net/http"

	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

// CustomResponder returns a custom response.
type CustomResponder struct {
	Message string `json:"message"`
}

func (c CustomResponder) ServeHTTP(w http.ResponseWriter, _ *http.Request, _ caddyhttp.Handler) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(c.Message))
	return err
}
