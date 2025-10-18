package tweet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	recentCountsPath = "/2/tweets/counts/recent"
	allCountsPath    = "/2/tweets/counts/all"
)

// Service coordinates tweet count operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the provided client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("recentcount: nil client")
	}
	return &Service{client: c}
}

// GetRecentCount fetches aggregate tweet counts for the recent endpoint.
func (s *Service) GetRecentCount(ctx context.Context, query, granularity string, params map[string]string) (CountResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return CountResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("recentcount: nil service")
	}

	qp := map[string]string{
		"query":       query,
		"granularity": granularity,
	}
	for key, value := range params {
		if key == "query" || key == "granularity" {
			continue
		}
		qp[key] = value
	}

	resp, err := s.client.Get(ctx, recentCountsPath, qp)
	if err != nil {
		return CountResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return CountResponse{}, rateLimits, err
	}

	var payload CountResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return CountResponse{}, rateLimits, fmt.Errorf("recentcount: decode response: %w", err)
	}

	return payload, rateLimits, nil
}

// GetAllCount queries the historical full-archive endpoint.
func (s *Service) GetAllCount(ctx context.Context, query, granularity string, params map[string]string) (CountResponse, client.RateLimitSnapshot, error) {
	qp := map[string]string{
		"query":       query,
		"granularity": granularity,
	}
	for key, value := range params {
		if key == "query" || key == "granularity" {
			continue
		}
		qp[key] = value
	}

	resp, err := s.client.Get(ctx, allCountsPath, qp)
	if err != nil {
		return CountResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return CountResponse{}, rateLimits, err
	}

	var payload CountResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return CountResponse{}, rateLimits, fmt.Errorf("recentcount: decode response: %w", err)
	}

	return payload, rateLimits, nil
}
