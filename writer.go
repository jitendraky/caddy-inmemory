package inmemory

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
)

type interceptingResponseWriter struct {
	io.Writer
	http.ResponseWriter
	Request *http.Request
}

func (w *interceptingResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *interceptingResponseWriter) Write(b []byte) (int, error) {
	w.Writer.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *interceptingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("not a Hijacker")
}

func (w *interceptingResponseWriter) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	} else {
		panic("not a Flusher")
	}
}

func (w *interceptingResponseWriter) CloseNotify() <-chan bool {
	if cn, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	panic("not a CloseNotifier")
}
