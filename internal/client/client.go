// Package client provides a configurable Twitter API HTTP client used by higher
// level service wrappers.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://api.twitter.com/"
	defaultUserAgent = "CERN-LineMode/2.15 libwww/2.17b3"
	defaultTimeout   = 60 * time.Second

	defaultRetryBase    = 500 * time.Millisecond
	defaultRetryWaitMax = 30 * time.Second
)

// Config describes how to construct a Client.
type Config struct {
	BaseURL     string
	BearerToken string
	UserAgent   string
	HTTPClient  *http.Client
	Timeout     time.Duration

	// Retry is the number of additional attempts made after a transient
	// failure (network error, HTTP 429, or a retryable 5xx). Zero disables
	// retries.
	Retry int

	// RetryWaitMax caps how long a single retry wait can last. Rate-limit
	// resets further away than this are truncated. Defaults to 30s.
	RetryWaitMax time.Duration
}

// Client wraps HTTP concerns for talking to the Twitter v2 API.
type Client struct {
	httpClient   *http.Client
	baseURL      *url.URL
	bearerToken  string
	userAgent    string
	retry        int
	retryBase    time.Duration
	retryWaitMax time.Duration
	logf         func(format string, args ...any)
}

// New constructs a Client using the supplied configuration taking sensible defaults
// from environment variables when values are omitted.
func New(cfg Config) (*Client, error) {
	baseURL, err := resolveBaseURL(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	bearer := strings.TrimSpace(cfg.BearerToken)
	if bearer == "" {
		bearer = strings.TrimSpace(os.Getenv("BEARER_TOKEN"))
	}

	userAgent := strings.TrimSpace(cfg.UserAgent)
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = defaultTimeout
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	retry := cfg.Retry
	if retry < 0 {
		retry = 0
	}

	retryWaitMax := cfg.RetryWaitMax
	if retryWaitMax <= 0 {
		retryWaitMax = defaultRetryWaitMax
	}

	return &Client{
		httpClient:   httpClient,
		baseURL:      baseURL,
		bearerToken:  bearer,
		userAgent:    userAgent,
		retry:        retry,
		retryBase:    defaultRetryBase,
		retryWaitMax: retryWaitMax,
		logf: func(format string, args ...any) {
			fmt.Fprintf(os.Stderr, format, args...)
		},
	}, nil
}

func resolveBaseURL(raw string) (*url.URL, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		trimmed = defaultBaseURL
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, fmt.Errorf("client: parse base url: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("client: invalid base url %q", trimmed)
	}

	if !strings.HasSuffix(parsed.Path, "/") {
		parsed.Path += "/"
	}

	return parsed, nil
}

// NewRequest builds an *http.Request using the configured base URL.
func (c *Client) NewRequest(ctx context.Context, method, path string, query map[string]string, body any) (*http.Request, error) {
	if c == nil {
		return nil, errors.New("client: nil Client")
	}
	if method == "" {
		return nil, errors.New("client: method must be provided")
	}

	fullURL, err := c.resolve(path)
	if err != nil {
		return nil, err
	}

	var reader io.Reader
	if body != nil {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		if err = encoder.Encode(body); err != nil {
			return nil, fmt.Errorf("client: encode body: %w", err)
		}
		reader = buffer
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("client: create request: %w", err)
	}

	if query != nil {
		values := req.URL.Query()
		for k, v := range query {
			values.Add(k, v)
		}
		req.URL.RawQuery = values.Encode()
	}

	c.decorateHeaders(req)
	return req, nil
}

// Get issues a GET request against the supplied path.
func (c *Client) Get(ctx context.Context, path string, query map[string]string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post issues a POST request against the supplied path with a JSON payload.
func (c *Client) Post(ctx context.Context, path string, body any, query map[string]string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, path, query, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Delete issues a DELETE request against the supplied path.
func (c *Client) Delete(ctx context.Context, path string, query map[string]string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodDelete, path, query, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Do forwards the request to the underlying http.Client while adding headers.
// Transient failures (network errors, HTTP 429, and retryable 5xx responses)
// are retried up to the configured number of attempts with exponential backoff,
// honoring Retry-After and x-rate-limit-reset headers when present.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil {
		return nil, errors.New("client: nil Client")
	}
	if req == nil {
		return nil, errors.New("client: nil request")
	}

	var lastErr error
	for attempt := 0; ; attempt++ {
		attemptReq, err := cloneRequest(req, attempt)
		if err != nil {
			if lastErr != nil {
				return nil, lastErr
			}
			return nil, err
		}

		resp, err := c.httpClient.Do(attemptReq)

		if attempt >= c.retry || !c.shouldRetry(req, resp, err) {
			return resp, err
		}

		wait := c.retryWait(attempt, resp)
		if err != nil {
			lastErr = err
			c.logf("ctw: request failed (%v); retrying in %s (attempt %d/%d)\n", err, wait.Round(time.Millisecond), attempt+1, c.retry)
		} else {
			lastErr = fmt.Errorf("client: transient status %d", resp.StatusCode)
			c.logf("ctw: got HTTP %d; retrying in %s (attempt %d/%d)\n", resp.StatusCode, wait.Round(time.Millisecond), attempt+1, c.retry)
			drainAndClose(resp.Body)
		}

		timer := time.NewTimer(wait)
		select {
		case <-req.Context().Done():
			timer.Stop()
			return nil, req.Context().Err()
		case <-timer.C:
		}
	}
}

// cloneRequest returns the request to use for the given attempt. The original
// request is used as-is for the first attempt; retries need a fresh body.
func cloneRequest(req *http.Request, attempt int) (*http.Request, error) {
	if attempt == 0 {
		return req, nil
	}
	clone := req.Clone(req.Context())
	if req.Body != nil {
		if req.GetBody == nil {
			return nil, errors.New("client: request body cannot be replayed for retry")
		}
		body, err := req.GetBody()
		if err != nil {
			return nil, fmt.Errorf("client: replay request body: %w", err)
		}
		clone.Body = body
	}
	return clone, nil
}

// shouldRetry reports whether a request may be safely attempted again.
func (c *Client) shouldRetry(req *http.Request, resp *http.Response, err error) bool {
	if req.Context().Err() != nil {
		return false
	}

	idempotent := req.Method == http.MethodGet || req.Method == http.MethodHead || req.Method == http.MethodDelete || req.Method == http.MethodPut

	if err != nil {
		// Network-level failure: only replay requests that are safe to repeat.
		return idempotent
	}
	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusServiceUnavailable:
		// Rate limited or unavailable: the request was not processed.
		return true
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
		// The server may have processed the request; only replay idempotent ones.
		return idempotent
	default:
		return false
	}
}

// retryWait computes how long to wait before the next attempt, preferring
// server-provided hints (Retry-After, x-rate-limit-reset) over exponential
// backoff, always bounded by retryWaitMax.
func (c *Client) retryWait(attempt int, resp *http.Response) time.Duration {
	base := c.retryBase
	if base <= 0 {
		base = defaultRetryBase
	}
	wait := base << uint(attempt)

	if resp != nil {
		if after := strings.TrimSpace(resp.Header.Get("Retry-After")); after != "" {
			if seconds, err := strconv.Atoi(after); err == nil && seconds > 0 {
				wait = time.Duration(seconds) * time.Second
			}
		} else if resp.StatusCode == http.StatusTooManyRequests {
			if reset := strings.TrimSpace(resp.Header.Get("x-rate-limit-reset")); reset != "" {
				if epoch, err := strconv.ParseInt(reset, 10, 64); err == nil {
					if until := time.Until(time.Unix(epoch, 0)); until > 0 {
						wait = until
					}
				}
			}
		}
	}

	wait = min(wait, c.retryWaitMax)
	if wait <= 0 {
		wait = base
	}
	return wait
}

func drainAndClose(body io.ReadCloser) {
	if body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 4096))
	_ = body.Close()
}

func (c *Client) resolve(path string) (*url.URL, error) {
	if c == nil {
		return nil, errors.New("client: nil Client")
	}
	if strings.TrimSpace(path) == "" {
		return c.baseURL, nil
	}

	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("client: parse path %q: %w", path, err)
	}
	return c.baseURL.ResolveReference(rel), nil
}

func (c *Client) decorateHeaders(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token := strings.TrimSpace(c.bearerToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if agent := strings.TrimSpace(c.userAgent); agent != "" {
		req.Header.Set("User-Agent", agent)
	}
}

// APIError represents one or more Twitter API errors.
type APIError struct {
	StatusCode int
	Errors     []Error
}

func (e APIError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("twitter api error: status %d", e.StatusCode)
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("twitter api error: status %d", e.StatusCode))
	for _, apiErr := range e.Errors {
		builder.WriteString(fmt.Sprintf("; code=%d message=%s", apiErr.Code, apiErr.Message))
	}
	return builder.String()
}

// Error captures a single entry in the Twitter error payload.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Errors []Error `json:"errors"`
}

// CheckResponse inspects the response status code and returns an APIError when
// the call did not succeed.
func CheckResponse(resp *http.Response) error {
	if resp == nil {
		return errors.New("client: nil response")
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer SafeClose(resp.Body)
	var payload errorResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("client: decode error response: %w", err)
	}

	return APIError{StatusCode: resp.StatusCode, Errors: payload.Errors}
}

// RateLimitSnapshot captures the rate limit headers for a response.
type RateLimitSnapshot struct {
	Limit     int
	Remaining int
	Reset     int
}

// ParseRateLimits extracts rate limiting information from the response headers.
func ParseRateLimits(resp *http.Response) RateLimitSnapshot {
	if resp == nil {
		return RateLimitSnapshot{Limit: -1, Remaining: -1, Reset: -1}
	}

	parseHeader := func(key string) int {
		value := strings.TrimSpace(resp.Header.Get(key))
		if value == "" {
			return -1
		}
		n, err := strconv.Atoi(value)
		if err != nil {
			return -1
		}
		return n
	}

	return RateLimitSnapshot{
		Limit:     parseHeader("x-rate-limit-limit"),
		Remaining: parseHeader("x-rate-limit-remaining"),
		Reset:     parseHeader("x-rate-limit-reset"),
	}
}

// SafeClose closes a response body while ignoring errors.
func SafeClose(closer io.ReadCloser) {
	if closer == nil {
		return
	}
	_ = closer.Close()
}
