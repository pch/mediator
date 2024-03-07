package internal

import (
	"log/slog"
	"net/http"
)

func NewImageHandler(config *Config, handler http.Handler) http.Handler {
	signatureMiddleware := NewSignatureMiddleware(config.SecretKey, handler)
	sourceMiddleware := NewSourceMiddleware(config, signatureMiddleware)
	cacheMiddleware := NewCacheMiddleware(config.CacheControl, sourceMiddleware)
	loggingMiddleware := NewLoggingMiddleware(cacheMiddleware)

	return loggingMiddleware
}

func NewHandler(config *Config) *http.ServeMux {
	transformHandler := NewImageHandler(config, NewImageTransformHandler(config))
	proxyHandler := NewImageHandler(config, NewProxyHandler(config))
	defaultRouteHandler := NewLoggingMiddleware(NewDefaultRouteHandler())

	pathPrefix := ensureValidPathPrefixFormat(config.PathPrefix)
	slog.Debug("Path prefix", "prefix", pathPrefix)

	mux := http.NewServeMux()
	mux.Handle(pathPrefix+"/image/transform/{source}/{path...}", transformHandler)
	mux.Handle(pathPrefix+"/proxy/{source}/{path...}", proxyHandler)
	mux.Handle("/", defaultRouteHandler)

	return mux
}
