package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func newTestClient(t *testing.T, baseURL string, retry int) *Client {
	t.Helper()
	c, err := New(Config{
		BaseURL:     baseURL,
		BearerToken: "test-token",
		Retry:       retry,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	c.retryBase = 5 * time.Millisecond
	c.logf = func(string, ...any) {}
	return c
}

func TestDoRetriesOn429(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 3)
	resp, err := c.Get(context.Background(), "2/tweets", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if got := atomic.LoadInt32(&calls); got != 2 {
		t.Fatalf("calls = %d, want 2", got)
	}
}

func TestDoRetriesPostOn503WithBodyReplay(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 3)
	resp, err := c.Post(context.Background(), "2/tweets", map[string]string{"text": "hello"}, nil)
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	defer SafeClose(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", resp.StatusCode)
	}
	if got := atomic.LoadInt32(&calls); got != 3 {
		t.Fatalf("calls = %d, want 3", got)
	}
}

func TestDoDoesNotRetryPostOn500(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 3)
	resp, err := c.Post(context.Background(), "2/tweets", map[string]string{"text": "hello"}, nil)
	if err != nil {
		t.Fatalf("Post: %v", err)
	}
	defer SafeClose(resp.Body)

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("calls = %d, want 1 (POST must not be replayed on 500)", got)
	}
}

func TestDoRetriesGetOn500(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 3)
	resp, err := c.Get(context.Background(), "2/tweets", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

func TestDoNoRetryWhenDisabled(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 0)
	resp, err := c.Get(context.Background(), "2/tweets", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer SafeClose(resp.Body)

	if got := atomic.LoadInt32(&calls); got != 1 {
		t.Fatalf("calls = %d, want 1", got)
	}
	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want 429 passed through", resp.StatusCode)
	}
}

func TestDoHonorsRetryAfterHeaderBounded(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 1)
	c.retryWaitMax = 50 * time.Millisecond // bound the 1s hint for the test

	start := time.Now()
	resp, err := c.Get(context.Background(), "2/tweets", nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("elapsed = %s, want wait capped by retryWaitMax", elapsed)
	}
}

func TestDoStopsOnContextCancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	c := newTestClient(t, server.URL, 3)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	_, err := c.Get(ctx, "2/tweets", nil)
	if err == nil {
		t.Fatal("expected context error, got nil")
	}
}
