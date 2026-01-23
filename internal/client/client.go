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
)

// Config describes how to construct a Client.
type Config struct {
	BaseURL     string
	BearerToken string
	UserAgent   string
	HTTPClient  *http.Client
	Timeout     time.Duration
}

// Client wraps HTTP concerns for talking to the Twitter v2 API.
type Client struct {
	httpClient  *http.Client
	baseURL     *url.URL
	bearerToken string
	userAgent   string
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

	return &Client{
		httpClient:  httpClient,
		baseURL:     baseURL,
		bearerToken: bearer,
		userAgent:   userAgent,
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
func (c *Client) NewRequest(ctx context.Context, method, path string, query map[string]string, body interface{}) (*http.Request, error) {
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
func (c *Client) Post(ctx context.Context, path string, body interface{}, query map[string]string) (*http.Response, error) {
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
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c == nil {
		return nil, errors.New("client: nil Client")
	}
	if req == nil {
		return nil, errors.New("client: nil request")
	}
	return c.httpClient.Do(req)
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
