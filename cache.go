package inmemory

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/allegro/bigcache"
)

type cache struct {
	responseCache *bigcache.BigCache
}

func newCache() *cache {

	config := bigcache.DefaultConfig(time.Second * 15)

	// @TODO(serafin): handle error
	bigCache, err := bigcache.NewBigCache(config)

	fmt.Println(err)

	return &cache{
		responseCache: bigCache,
	}
}

func (c cache) variantKey(writer http.ResponseWriter, r *http.Request) string {
	var values []string

	headers := writer.Header()

	fmt.Printf("Response headers: %v\n", headers)
	fmt.Printf("Request headers: %v", r.Header)

	values = append(values, fmt.Sprintf("URL:%s", r.URL.String()))

	// @TODO: normalize/clean headers to prevent DDoSing cache
	for _, vary := range headers["Vary"] {
		if requestHeader, exists := r.Header[vary]; exists {
			values = append(values, fmt.Sprintf("%s:%v", vary, requestHeader))
		}
	}

	return strings.Join(values, ";;")
}

func (c cache) knownVariants(r *http.Request) []variant {

	if bytes, err := c.responseCache.Get(c.requestCacheKey(r)); err == nil {
		return decodeVariants(bytes)
	}

	return []variant{}
}

func (c cache) requestCacheKey(r *http.Request) string {
	return r.URL.String()
}

func (c cache) cacheResponse(code int, buffer *bytes.Buffer, writer http.ResponseWriter, r *http.Request) {
	requestKey := c.requestCacheKey(r)
	variantKey := c.variantKey(writer, r)

	var variants []variant

	if buf, err := c.responseCache.Get(requestKey); err == nil {
		variants = decodeVariants(buf)

		if !c.hasVariant(variants, variantKey) {
			variants = append(variants, variant{
				CacheKey: variantKey,
			})

			c.responseCache.Set(requestKey, encodeVariants(variants))
		}
	} else {
		variants = append(variants, variant{
			CacheKey: variantKey,
		})

		c.responseCache.Set(requestKey, encodeVariants(variants))
	}

	c.responseCache.Set(variantKey, encodeCachedResponse(cachedResponse{
		Headers: writer.Header(),
		Bytes:   buffer.Bytes(),
		Code:    code,
	}))
}

func (c cache) getFromCache(r *http.Request) (cachedResponse, error) {
	requestKey := c.requestCacheKey(r)

	if buf, err := c.responseCache.Get(requestKey); err == nil {
		variants := decodeVariants(buf)

		if len(variants) > 0 {
			if buf2, err2 := c.responseCache.Get(variants[0].CacheKey); err2 == nil {
				return decodeCachedResponse(buf2), nil
			}
		}
	}

	return cachedResponse{}, fmt.Errorf("Cache not found")
}

func (c cache) hasVariant(variants []variant, key string) bool {
	for _, v := range variants {
		if v.CacheKey == key {
			return true
		}
	}

	return false
}
