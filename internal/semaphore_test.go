package internal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphoreLimitsConcurrentTransforms(t *testing.T) {
	const maxConcurrent = 2
	const totalRequests = 6

	var running atomic.Int32
	var peaked atomic.Int32
	gate := make(chan struct{})

	testImage := makePNG(t, 40, 20)

	// Upstream image server that blocks until we release the gate.
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cur := running.Add(1)
		defer running.Add(-1)

		// Track peak concurrency.
		for {
			old := peaked.Load()
			if cur <= old || peaked.CompareAndSwap(old, cur) {
				break
			}
		}

		<-gate

		w.Header().Set("Content-Type", "image/png")
		w.Write(testImage)
	}))
	defer upstream.Close()

	config := &Config{
		DownloadMaxSize:         50 * MB,
		DownloadTimeout:         5 * time.Second,
		MaxConcurrentTransforms: maxConcurrent,
		CacheControl:            defaultCacheControl,
		SecretKey:               "",
	}

	handler := NewImageTransformHandler(config)

	var wg sync.WaitGroup
	responses := make([]*httptest.ResponseRecorder, totalRequests)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		responses[i] = httptest.NewRecorder()

		go func(idx int) {
			defer wg.Done()

			req := httptest.NewRequest("GET", "/image/transform/images/photo.png?op=fit&w=10&h=10", nil)
			ctx := setImageSource(req.Context(), &ImageSource{URL: upstream.URL + "/photo.png"})
			req = req.WithContext(ctx)

			handler.ServeHTTP(responses[idx], req)
		}(i)
	}

	// Wait for the semaphore to fill up.
	deadline := time.After(5 * time.Second)
	for {
		if running.Load() >= int32(maxConcurrent) {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for %d goroutines to start, got %d", maxConcurrent, running.Load())
		default:
			time.Sleep(5 * time.Millisecond)
		}
	}

	// Give a moment for any extra goroutines to sneak through.
	time.Sleep(50 * time.Millisecond)

	if peak := peaked.Load(); peak > int32(maxConcurrent) {
		t.Fatalf("peak concurrency = %d, want <= %d", peak, maxConcurrent)
	}

	// Release all blocked handlers.
	close(gate)
	wg.Wait()
}

func TestSemaphoreRespectsRequestCancellation(t *testing.T) {
	gate := make(chan struct{})
	testImage := makePNG(t, 40, 20)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-gate
		w.Header().Set("Content-Type", "image/png")
		w.Write(testImage)
	}))
	defer upstream.Close()

	config := &Config{
		DownloadMaxSize:         50 * MB,
		DownloadTimeout:         5 * time.Second,
		MaxConcurrentTransforms: 1,
		CacheControl:            defaultCacheControl,
		SecretKey:               "",
	}

	handler := NewImageTransformHandler(config)

	// Fill the semaphore with one blocking request.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		req := httptest.NewRequest("GET", "/image/transform/images/photo.png?op=fit&w=10&h=10", nil)
		ctx := setImageSource(req.Context(), &ImageSource{URL: upstream.URL + "/photo.png"})
		req = req.WithContext(ctx)
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}()

	// Wait for the first request to acquire the semaphore.
	time.Sleep(50 * time.Millisecond)

	// Send a second request with an already-cancelled context.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest("GET", "/image/transform/images/photo.png?op=fit&w=10&h=10", nil)
	req = req.WithContext(setImageSource(ctx, &ImageSource{URL: upstream.URL + "/photo.png"}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}

	close(gate)
	wg.Wait()
}
