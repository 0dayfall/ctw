package publish

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	createTweetPath = "/2/tweets"
	deleteTweetFmt  = "/2/tweets/%s"
)

// Service coordinates tweet create/delete operations.
// Prefer constructing a single instance via NewService and reusing it across commands.
type Service struct {
	client *client.Client
}

// NewService constructs a publish Service backed by the provided client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("publish: nil client")
	}
	return &Service{client: c}
}

// CreateTweet issues a POST /2/tweets call with the provided payload.
func (s *Service) CreateTweet(ctx context.Context, req CreateTweetRequest) (CreateTweetResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return CreateTweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("publish: nil service")
	}

	resp, err := s.client.Post(ctx, createTweetPath, req, nil)
	if err != nil {
		return CreateTweetResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return CreateTweetResponse{}, rateLimits, err
	}

	var payload CreateTweetResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return CreateTweetResponse{}, rateLimits, fmt.Errorf("publish: decode create tweet response: %w", err)
	}

	return payload, rateLimits, nil
}

// DeleteTweet removes an existing tweet by ID using DELETE /2/tweets/:id.
func (s *Service) DeleteTweet(ctx context.Context, tweetID string) (DeleteTweetResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return DeleteTweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("publish: nil service")
	}

	path := fmt.Sprintf(deleteTweetFmt, tweetID)
	resp, err := s.client.Delete(ctx, path, nil)
	if err != nil {
		return DeleteTweetResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return DeleteTweetResponse{}, rateLimits, err
	}

	var payload DeleteTweetResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return DeleteTweetResponse{}, rateLimits, fmt.Errorf("publish: decode delete tweet response: %w", err)
	}

	return payload, rateLimits, nil
}
