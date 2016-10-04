package inmemory

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

type cacheHandler struct {
	Next   httpserver.Handler
	Config config
}

func (h cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	b := bytes.NewBuffer(nil)
	rw := &interceptingResponseWriter{Writer: b, ResponseWriter: w, Request: r}

	code, err := h.Next.ServeHTTP(rw, r)

	fmt.Printf("Intercepted %d bytes with code %d, error: %v, request: %v\n", len(b.Bytes()), code, err, r)
	fmt.Printf("Response headers are: %v", w.Header())

	return code, err
}
