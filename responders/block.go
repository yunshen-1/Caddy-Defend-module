package responders

import (
	"net/http"
)

// BlockResponder blocks the request with a 403 Forbidden response.
type BlockResponder struct{}

func (b BlockResponder) Respond(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusForbidden)
	_, err := w.Write([]byte("Access denied"))
	return err
}
