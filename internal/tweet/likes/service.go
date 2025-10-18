package likes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	likePathTemplate    = "/2/users/%s/likes"
	unlikePathTemplate  = "/2/users/%s/likes/%s"
	likedTweetsTemplate = "/2/users/%s/liked_tweets"
)

// Service wraps Twitter like endpoints.
type Service struct {
	client *client.Client
}

// NewService constructs a Service using the provided client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("likes: nil client")
	}
	return &Service{client: c}
}

// LikeTweet creates a like relationship for the given tweet on behalf of userID.
func (s *Service) LikeTweet(ctx context.Context, userID, tweetID string) (RelationshipResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return RelationshipResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("likes: nil service")
	}

	path := fmt.Sprintf(likePathTemplate, userID)
	payload := map[string]string{"tweet_id": tweetID}
	resp, err := s.client.Post(ctx, path, payload, nil)
	if err != nil {
		return RelationshipResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RelationshipResponse{}, rateLimits, err
	}

	var relationship RelationshipResponse
	if err := json.NewDecoder(resp.Body).Decode(&relationship); err != nil {
		return RelationshipResponse{}, rateLimits, fmt.Errorf("likes: decode like response: %w", err)
	}

	return relationship, rateLimits, nil
}

// UnlikeTweet removes a like relationship for the given tweet on behalf of userID.
func (s *Service) UnlikeTweet(ctx context.Context, userID, tweetID string) (RelationshipResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return RelationshipResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("likes: nil service")
	}

	path := fmt.Sprintf(unlikePathTemplate, userID, tweetID)
	resp, err := s.client.Delete(ctx, path, nil)
	if err != nil {
		return RelationshipResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return RelationshipResponse{}, rateLimits, err
	}

	var relationship RelationshipResponse
	if err := json.NewDecoder(resp.Body).Decode(&relationship); err != nil {
		return RelationshipResponse{}, rateLimits, fmt.Errorf("likes: decode unlike response: %w", err)
	}

	return relationship, rateLimits, nil
}

// ListLikedTweets retrieves liked tweets for userID applying optional query params (pagination, expansions, etc.).
func (s *Service) ListLikedTweets(ctx context.Context, userID string, params map[string]string) (LikedTweetsResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return LikedTweetsResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("likes: nil service")
	}

	path := fmt.Sprintf(likedTweetsTemplate, userID)
	resp, err := s.client.Get(ctx, path, params)
	if err != nil {
		return LikedTweetsResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return LikedTweetsResponse{}, rateLimits, err
	}

	var payload LikedTweetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return LikedTweetsResponse{}, rateLimits, fmt.Errorf("likes: decode liked tweets response: %w", err)
	}

	return payload, rateLimits, nil
}
