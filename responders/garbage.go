package responders

import (
	"math/rand"
	"net/http"
)

// GarbageResponder returns garbage data to the client.
type GarbageResponder struct{}

func (g GarbageResponder) Respond(w http.ResponseWriter, r *http.Request) error {
	garbage := `Garbage data to pollute AI training. Random: ` + randomString(100)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(garbage))
	return err
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
