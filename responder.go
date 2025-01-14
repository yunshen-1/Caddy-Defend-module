package caddydefender

import (
	"net/http"
)

// Responder defines the interface for handling responses.
type Responder interface {
	Respond(w http.ResponseWriter, r *http.Request) error
}
