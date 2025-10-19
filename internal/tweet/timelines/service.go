package timelines

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	userTweetsPathFormat    = "/2/users/%s/tweets"
	userMentionsPathFormat  = "/2/users/%s/mentions"
	reverseChronoPathFormat = "/2/users/%s/timelines/reverse_chronological"
)

// Service coordinates Twitter timeline operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the supplied client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("timelines: nil client")
	}
	return &Service{client: c}
}

// GetUserTweets fetches tweets posted by the specified user.
func (s *Service) GetUserTweets(ctx context.Context, userID string, params map[string]string) (TimelineResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return TimelineResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("timelines: nil service")
	}
	if userID == "" {
		return TimelineResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("timelines: user id is required")
	}

	path := fmt.Sprintf(userTweetsPathFormat, userID)
	return s.fetch(ctx, path, params)
}

// GetUserMentions fetches tweets that mention the specified user.
func (s *Service) GetUserMentions(ctx context.Context, userID string, params map[string]string) (TimelineResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return TimelineResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("timelines: nil service")
	}
	if userID == "" {
		return TimelineResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("timelines: user id is required")
	}

	path := fmt.Sprintf(userMentionsPathFormat, userID)
	return s.fetch(ctx, path, params)
}

// GetReverseChronological fetches the reverse chronological home timeline for the authenticated user.
func (s *Service) GetReverseChronological(ctx context.Context, userID string, params map[string]string) (TimelineResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return TimelineResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("timelines: nil service")
	}
	if userID == "" {
		return TimelineResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("timelines: user id is required")
	}

	path := fmt.Sprintf(reverseChronoPathFormat, userID)
	return s.fetch(ctx, path, params)
}

func (s *Service) fetch(ctx context.Context, path string, params map[string]string) (TimelineResponse, client.RateLimitSnapshot, error) {
	resp, err := s.client.Get(ctx, path, params)
	if err != nil {
		return TimelineResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return TimelineResponse{}, rateLimits, err
	}

	var payload TimelineResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return TimelineResponse{}, rateLimits, fmt.Errorf("timelines: decode response: %w", err)
	}

	return payload, rateLimits, nil
}
