package retweets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	retweetPathFormat    = "/2/users/%s/retweets"
	unretweetPathFormat  = "/2/users/%s/retweets/%s"
	retweetersPathFormat = "/2/tweets/%s/retweeted_by"
)

// Service coordinates Twitter retweet operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the supplied client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("retweets: nil client")
	}
	return &Service{client: c}
}

// Retweet creates a retweet on behalf of the specified user.
func (s *Service) Retweet(ctx context.Context, userID, tweetID string) (RetweetResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return RetweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: nil service")
	}
	if userID == "" {
		return RetweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: user id is required")
	}
	if tweetID == "" {
		return RetweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: tweet id is required")
	}

	path := fmt.Sprintf(retweetPathFormat, userID)
	payload := RetweetRequest{TweetID: tweetID}

	resp, err := s.client.Post(ctx, path, payload, nil)
	if err != nil {
		return RetweetResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RetweetResponse{}, rateLimits, err
	}

	var result RetweetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RetweetResponse{}, rateLimits, fmt.Errorf("retweets: decode response: %w", err)
	}

	return result, rateLimits, nil
}

// Unretweet removes a retweet on behalf of the specified user.
func (s *Service) Unretweet(ctx context.Context, userID, tweetID string) (RetweetResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return RetweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: nil service")
	}
	if userID == "" {
		return RetweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: user id is required")
	}
	if tweetID == "" {
		return RetweetResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: tweet id is required")
	}

	path := fmt.Sprintf(unretweetPathFormat, userID, tweetID)

	resp, err := s.client.Delete(ctx, path, nil)
	if err != nil {
		return RetweetResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RetweetResponse{}, rateLimits, err
	}

	var result RetweetResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RetweetResponse{}, rateLimits, fmt.Errorf("retweets: decode response: %w", err)
	}

	return result, rateLimits, nil
}

// ListRetweeters fetches users who retweeted the specified tweet.
func (s *Service) ListRetweeters(ctx context.Context, tweetID string, params map[string]string) (RetweetersResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return RetweetersResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: nil service")
	}
	if tweetID == "" {
		return RetweetersResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("retweets: tweet id is required")
	}

	path := fmt.Sprintf(retweetersPathFormat, tweetID)

	resp, err := s.client.Get(ctx, path, params)
	if err != nil {
		return RetweetersResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RetweetersResponse{}, rateLimits, err
	}

	var result RetweetersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return RetweetersResponse{}, rateLimits, fmt.Errorf("retweets: decode response: %w", err)
	}

	return result, rateLimits, nil
}
