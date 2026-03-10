package internal

import (
	"net"
	"net/http"
	"testing"
	"time"
)

func TestHttpServerStartReturnsErrorWhenPortInUse(t *testing.T) {
	busyListener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen busy port: %v", err)
	}
	defer busyListener.Close()

	port := busyListener.Addr().(*net.TCPAddr).Port

	cfg := &Config{
		HttpPort:         port,
		HttpIdleTimeout:  time.Second,
		HttpReadTimeout:  time.Second,
		HttpWriteTimeout: time.Second,
	}

	s := NewHttpServer(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	if err := s.Start(); err == nil {
		t.Fatalf("expected startup error for busy port")
	}
}

func TestHttpServerStartAndStop(t *testing.T) {
	cfg := &Config{
		HttpPort:         0,
		HttpIdleTimeout:  time.Second,
		HttpReadTimeout:  time.Second,
		HttpWriteTimeout: time.Second,
	}

	s := NewHttpServer(cfg, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	if err := s.Start(); err != nil {
		t.Fatalf("Start() error: %v", err)
	}

	// Should not panic.
	s.Stop()
}

func TestHttpServerStopWithoutStart(t *testing.T) {
	s := &HttpServer{}
	// Should not panic.
	s.Stop()
}
