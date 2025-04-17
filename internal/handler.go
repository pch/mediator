package internal

import (
	"log/slog"
	"net/http"
)

func NewSourcedMediaHandler(config *Config, handler http.Handler) http.Handler {
	signatureMiddleware := NewSignatureMiddleware(config.SecretKey, handler)
	imageSourceMiddleware := NewImageSourceMiddleware(config, signatureMiddleware)
	authMiddleware := NewAuthMiddleware(config, imageSourceMiddleware)
	loggingMiddleware := NewLoggingMiddleware(authMiddleware)

	return loggingMiddleware
}

func NewUnsourcedMediaHandler(config *Config, handler http.Handler) http.Handler {
	signatureMiddleware := NewSignatureMiddleware(config.SecretKey, handler)
	authMiddleware := NewAuthMiddleware(config, signatureMiddleware)
	loggingMiddleware := NewLoggingMiddleware(authMiddleware)

	return loggingMiddleware
}

func NewHandler(config *Config) *http.ServeMux {
	transformHandler := NewSourcedMediaHandler(config, NewImageTransformHandler(config))
	renderHandler := NewUnsourcedMediaHandler(config, NewRenderHandler(config))
	defaultRouteHandler := NewLoggingMiddleware(NewDefaultRouteHandler())

	pathPrefix := ensureValidPathPrefixFormat(config.PathPrefix)
	slog.Debug("Path prefix", "prefix", pathPrefix)

	mux := http.NewServeMux()
	mux.Handle(pathPrefix+"/image/transform/{source}/{path...}", transformHandler)
	mux.Handle(pathPrefix+"/render/{renderer}/{payloadBase64}", renderHandler)
	mux.Handle("/", defaultRouteHandler)

	return mux
}
