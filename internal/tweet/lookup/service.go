package lookup

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	tweetsLookupPath = "/2/tweets"
	tweetLookupPath  = "/2/tweets/%s"
)

// Service coordinates Twitter tweet lookup operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the supplied client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("lookup: nil client")
	}
	return &Service{client: c}
}

// GetTweet fetches a single tweet by ID with optional query parameters.
func (s *Service) GetTweet(ctx context.Context, tweetID string, params map[string]string) (TweetLookupResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return TweetLookupResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("lookup: nil service")
	}
	if strings.TrimSpace(tweetID) == "" {
		return TweetLookupResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("lookup: tweet id is required")
	}

	path := fmt.Sprintf(tweetLookupPath, tweetID)

	resp, err := s.client.Get(ctx, path, params)
	if err != nil {
		return TweetLookupResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return TweetLookupResponse{}, rateLimits, err
	}

	var result TweetLookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return TweetLookupResponse{}, rateLimits, fmt.Errorf("lookup: decode response: %w", err)
	}

	return result, rateLimits, nil
}

// GetTweets fetches multiple tweets by IDs with optional query parameters.
func (s *Service) GetTweets(ctx context.Context, tweetIDs []string, params map[string]string) (TweetLookupResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return TweetLookupResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("lookup: nil service")
	}
	if len(tweetIDs) == 0 {
		return TweetLookupResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("lookup: tweet ids are required")
	}

	queryParams := make(map[string]string)
	for k, v := range params {
		queryParams[k] = v
	}
	queryParams["ids"] = strings.Join(tweetIDs, ",")

	resp, err := s.client.Get(ctx, tweetsLookupPath, queryParams)
	if err != nil {
		return TweetLookupResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return TweetLookupResponse{}, rateLimits, err
	}

	var result TweetLookupResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return TweetLookupResponse{}, rateLimits, fmt.Errorf("lookup: decode response: %w", err)
	}

	return result, rateLimits, nil
}
