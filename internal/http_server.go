package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type HttpServer struct {
	config     *Config
	handler    http.Handler
	httpServer *http.Server
}

func NewHttpServer(config *Config, handler http.Handler) *HttpServer {
	return &HttpServer{
		handler: handler,
		config:  config,
	}
}

func (s *HttpServer) Start() {
	httpAddress := fmt.Sprintf(":%d", s.config.HttpPort)

	slog.Info("Server starting", "http", httpAddress)

	s.httpServer = s.newHttpServer(httpAddress)
	s.httpServer.Handler = s.handler

	go s.httpServer.ListenAndServe()

	slog.Info("Server started", "http", httpAddress)
}

func (s *HttpServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	defer slog.Info("Server stopped")

	slog.Info("Server stopping...")

	s.httpServer.Shutdown(ctx)
}

func (s *HttpServer) newHttpServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		IdleTimeout:  s.config.HttpIdleTimeout,
		ReadTimeout:  s.config.HttpReadTimeout,
		WriteTimeout: s.config.HttpWriteTimeout,
	}
}
