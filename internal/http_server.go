package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
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

func (s *HttpServer) Start() error {
	httpAddress := fmt.Sprintf(":%d", s.config.HttpPort)

	slog.Info("Server starting", "http", httpAddress)

	s.httpServer = s.newHttpServer(httpAddress)
	s.httpServer.Handler = s.handler

	listener, err := net.Listen("tcp", httpAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", httpAddress, err)
	}

	go func() {
		if err := s.httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "http", httpAddress, "error", err)
		}
	}()

	slog.Info("Server started", "http", httpAddress)

	return nil
}

func (s *HttpServer) Stop() {
	if s.httpServer == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	slog.Info("Server stopping...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		return
	}

	slog.Info("Server stopped")
}

func (s *HttpServer) newHttpServer(addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		IdleTimeout:  s.config.HttpIdleTimeout,
		ReadTimeout:  s.config.HttpReadTimeout,
		WriteTimeout: s.config.HttpWriteTimeout,
	}
}
