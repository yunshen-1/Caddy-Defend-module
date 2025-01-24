package responders

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"net/http"
)

// BlockResponder blocks the request with a 403 Forbidden response.
type BlockResponder struct{}

func (b BlockResponder) ServeHTTP(w http.ResponseWriter, _ *http.Request, _ caddyhttp.Handler) error {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write([]byte("Access denied"))
	return err
}
