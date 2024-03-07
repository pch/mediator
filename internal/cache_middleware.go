package internal

import "net/http"

type CacheMiddleware struct {
	cacheControl string
	next         http.Handler
}

func NewCacheMiddleware(cacheControl string, next http.Handler) *CacheMiddleware {
	return &CacheMiddleware{cacheControl, next}
}

func (h *CacheMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	newWriter := newResponseWriter(w)

	h.next.ServeHTTP(newWriter, r)

	if newWriter.statusCode == http.StatusOK {
		w.Header().Set("Cache-Control", h.cacheControl)
	}
}
