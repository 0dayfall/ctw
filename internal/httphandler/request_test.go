package httphandler

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGetRequest(t *testing.T) {
	assert := assert.New(t)

	testURL := "http://example.com"
	httpRequest := CreateGetRequest(testURL)

	assert.Equal(httpRequest.Method, "GET")
	assert.Equal(httpRequest.URL.String(), testURL)
	assert.Equal(httpRequest.Header.Get("Accept-Encoding"), "gzip")
}

func TestIsRateLimitReset429(t *testing.T) {
	// Create a mock HTTP response with status code 429 (Too Many Requests)
	// the header x-rate-limit-reset is set to 60 seconds
	// error code is 88 with message "rate limit exceeded"
	mockResponse := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     map[string][]string{},
		Body:       io.NopCloser(bytes.NewReader([]byte(`{ "errors": [ { "code": 88, "message": "Rate limit exceeded" } ] })`))),
	}
	mockResponse.Header.Add("x-rate-limit-reset", "60")

	isOK, rateLimitReset := IsRateLimitResetOK(mockResponse)

	require.Falsef(t, isOK, "Expected IsRateLimitOK to return false, got %v", isOK)
	require.EqualValuesf(t, 60, rateLimitReset, "Expected timeToReset to be 60, got %d", rateLimitReset)
}

func TestIsRateLimitResetOK(t *testing.T) {
	// Create a mock HTTP response with status code 200 (OK)
	mockResponse := &http.Response{
		StatusCode: http.StatusOK,
		Header:     map[string][]string{},
	}
	mockResponse.Header.Add("x-rate-limit-reset", "120")
	isOK, rateLimitReset := IsRateLimitResetOK(mockResponse)

	require.Truef(t, isOK, "Expected IsRateLimitOK to return true, got %v", isOK)
	require.EqualValuesf(t, 120, rateLimitReset, "Expected timeToReset to be 120, got %d", rateLimitReset)
}

func TestRateLimitLimit429(t *testing.T) {
	mockResponse := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     map[string][]string{},
		Body:       io.NopCloser(bytes.NewReader([]byte(`{ "errors": [ { "code": 88, "message": "Rate limit exceeded" } ] })`))),
	}
	mockResponse.Header.Add("x-rate-limit-limit", "60")
	limit := GetRateLimitLimit(mockResponse)
	require.EqualValuesf(t, 60, limit, "Expected timeToReset to be 60, got %d", limit)
}

func TestRateLimitRemaining429(t *testing.T) {
	mockResponse := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     map[string][]string{},
		Body:       io.NopCloser(bytes.NewReader([]byte(`{ "errors": [ { "code": 88, "message": "You have exceeded your API rate limit. Please wait before making more requests." } ] })`))),
	}
	mockResponse.Header.Add("x-rate-limit-remaining", "0")
	timeToReset := GetRateLimitRemaining(mockResponse)
	require.EqualValuesf(t, 0, timeToReset, "Expected timeToReset to be 60, got %d", timeToReset)
}
