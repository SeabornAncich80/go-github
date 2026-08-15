package main

import (
	"net/http"
	"sync"
)

// CachingTransport implements http.RoundTripper with ETag cache invalidation.
type CachingTransport struct {
	Base  http.RoundTripper
	cache sync.Map
	mu    sync.RWMutex
}

func (t *CachingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Handle mutation requests
	if isMutation(req.Method) {
		resp, err := t.Base.RoundTrip(req)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			t.invalidate(req.URL.String())
			// Invalidate Location/Content-Location headers
			for _, header := range []string{"Location", "Content-Location"} {
				if loc := resp.Header.Get(header); loc != "" {
					t.invalidate(loc)
				}
			}
		}
		return resp, err
	}

	// Standard GET/HEAD logic would go here
	return t.Base.RoundTrip(req)
}

func (t *CachingTransport) invalidate(url string) {
	t.cache.Delete(url)
}

func isMutation(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

func main() {
	// Implementation placeholder
}