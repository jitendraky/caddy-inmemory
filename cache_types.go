package inmemory

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type variant struct {
	CacheKey string
}

type cachedResponse struct {
	Headers map[string][]string
	Bytes   []byte
	Code    int
}

func (c cachedResponse) String() string {
	return fmt.Sprintf("CachedResponse{headers: %v, code: %d, content: %d}", c.Headers, c.Code, len(c.Bytes))
}

func encodeCachedResponse(response cachedResponse) []byte {
	var bytes bytes.Buffer
	encoder := gob.NewEncoder(&bytes)
	encoder.Encode(response)

	return bytes.Bytes()
}

func decodeCachedResponse(b []byte) cachedResponse {
	buffer := bytes.NewBuffer(b)

	var result cachedResponse
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(&result)

	return result
}

func encodeVariants(variants []variant) []byte {
	var bytes bytes.Buffer
	encoder := gob.NewEncoder(&bytes)
	encoder.Encode(variants)

	return bytes.Bytes()
}

func decodeVariants(b []byte) []variant {
	buffer := bytes.NewBuffer(b)

	var result []variant
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(&result)

	return result
}
