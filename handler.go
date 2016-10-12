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
	Cache  *cache
}

func (ch cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	if ch.isPurgeRequest(r) {
		return ch.handlePurgeRequest(w, r)
	}

	if ch.shouldPassThroughCache(r) {
		return ch.interceptAndCacheResponse(w, r)
	}

	return ch.Next.ServeHTTP(w, r)
}

func (ch cacheHandler) shouldPassThroughCache(r *http.Request) bool {
	// @TODO(serafin) check if URL is matching filters
	fmt.Printf("shouldPassThroughCache(%v)\n", r.URL)
	return true
}

func (ch cacheHandler) interceptAndCacheResponse(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Printf("interceptAndCacheResponse(%v)\n", r.URL)

	if resp, err := ch.getCachedResponse(r); err == nil {

		fmt.Printf("Response is: %v\n", resp)

		for name, values := range resp.Headers {

			fmt.Printf("name: %s, values: %s\n", name, values)

			for index, value := range values {
				if index == 0 {
					w.Header().Set(name, value)
				} else {
					w.Header().Add(name, value)
				}
			}
		}

		w.WriteHeader(resp.Code)
		w.Write(resp.Bytes)

		return resp.Code, nil
	}

	return ch.sendDownStreamAndCache(w, r)
}

func (ch cacheHandler) getCachedKey(r *http.Request) string {
	fmt.Printf("isCached(%v)\n", r.URL)

	return ch.Cache.requestCacheKey(r)
}

func (ch cacheHandler) sendDownStreamAndCache(w http.ResponseWriter, r *http.Request) (int, error) {

	fmt.Printf("sendDownStreamAndCache(%v)\n", r.URL)

	buffer := bytes.NewBuffer(nil)
	rw := &interceptingResponseWriter{Writer: buffer, ResponseWriter: w, Request: r}
	code, err := ch.Next.ServeHTTP(rw, r)

	if ch.shouldStoreInCache(rw, r, code, err) {
		go ch.cacheResponse(code, buffer, rw, r)
	}

	return code, err
}

func (ch cacheHandler) getCachedResponse(r *http.Request) (cachedResponse, error) {
	return ch.Cache.getFromCache(r)
}

func (ch cacheHandler) cacheResponse(code int, buffer *bytes.Buffer, writer http.ResponseWriter, r *http.Request) {
	ch.Cache.cacheResponse(code, buffer, writer, r)
}

func (ch cacheHandler) shouldStoreInCache(writer http.ResponseWriter, r *http.Request, code int, err error) bool {
	fmt.Printf("shouldStoreInCache(%v)\n", r.URL)

	headers := writer.Header()

	// Cache only responses containing ETags
	if headers.Get("ETag") == "" {
		fmt.Printf("ETag is empty!!!")
		return false
	}

	// Cannot cache session-aware responses
	if headers.Get("Set-Cookie") != "" {
		fmt.Printf("Response has cookies!!!")
		return false
	}

	// Cache only HTTP 200
	if code != http.StatusOK {
		fmt.Printf("Response is not ok!!!")
		return false
	}

	//@TODO check if URL matches / is baned

	return true
}

func (ch cacheHandler) handlePurgeRequest(w http.ResponseWriter, r *http.Request) (int, error) {
	return 0, nil
}

func (ch cacheHandler) isPurgeRequest(r *http.Request) bool {
	return false
}
