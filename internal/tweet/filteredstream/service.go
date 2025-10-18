package tweet

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	streamPath = "/2/tweets/search/stream"
	rulesPath  = "/2/tweets/search/stream/rules"
)

// Service coordinates filtered stream operations against the Twitter API.
type Service struct {
	client *client.Client
}

// NewService returns a Service backed by the provided client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("filteredstream: nil client")
	}
	return &Service{client: c}
}

// Stream fetches filtered stream results using the provided expansion fields.
func (s *Service) Stream(ctx context.Context, fields map[string]string) (StreamEnvelope, client.RateLimitSnapshot, error) {
	resp, err := s.client.Get(ctx, streamPath, fields)
	if err != nil {
		return StreamEnvelope{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return StreamEnvelope{}, rateLimits, err
	}

	var payload StreamEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return StreamEnvelope{}, rateLimits, fmt.Errorf("filteredstream: decode stream response: %w", err)
	}

	return payload, rateLimits, nil
}

// AddRule creates a new filtered stream rule.
func (s *Service) AddRule(ctx context.Context, cmd AddCommand, dryRun bool) (RulesResponse, client.RateLimitSnapshot, error) {
	query := map[string]string{}
	if dryRun {
		query["dry_run"] = "true"
	}

	resp, err := s.client.Post(ctx, rulesPath, cmd, query)
	if err != nil {
		return RulesResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RulesResponse{}, rateLimits, err
	}

	var payload RulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return RulesResponse{}, rateLimits, fmt.Errorf("filteredstream: decode add rule response: %w", err)
	}

	return payload, rateLimits, nil
}

// GetRules retrieves the current filtered stream rules.
func (s *Service) GetRules(ctx context.Context) (RulesResponse, client.RateLimitSnapshot, error) {
	resp, err := s.client.Get(ctx, rulesPath, nil)
	if err != nil {
		return RulesResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RulesResponse{}, rateLimits, err
	}

	var payload RulesResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return RulesResponse{}, rateLimits, fmt.Errorf("filteredstream: decode get rules response: %w", err)
	}

	return payload, rateLimits, nil
}
