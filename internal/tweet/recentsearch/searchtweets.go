package tweet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const recentSearchPath = "/2/tweets/search/recent"

// Service exposes helpers for the recent search endpoints.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the provided client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("recentsearch: nil client")
	}
	return &Service{client: c}
}

// SearchRecent queries the recent search endpoint with optional query parameters.
// The query string is always applied while the params map can be used for
// pagination or additional expansions.
func (s *Service) SearchRecent(ctx context.Context, query string, params map[string]string) (SearchRecentResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return SearchRecentResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("recentsearch: nil service")
	}

	qp := map[string]string{"query": query}
	for key, value := range params {
		if key == "query" {
			continue
		}
		qp[key] = value
	}

	resp, err := s.client.Get(ctx, recentSearchPath, qp)
	if err != nil {
		return SearchRecentResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return SearchRecentResponse{}, rateLimits, err
	}

	var payload SearchRecentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return SearchRecentResponse{}, rateLimits, fmt.Errorf("recentsearch: decode response: %w", err)
	}

	return payload, rateLimits, nil
}

// SearchRecentNextToken is a convenience wrapper that applies a pagination token.
func (s *Service) SearchRecentNextToken(ctx context.Context, query, token string) (SearchRecentResponse, client.RateLimitSnapshot, error) {
	params := map[string]string{"pagination_token": token}
	return s.SearchRecent(ctx, query, params)
}
