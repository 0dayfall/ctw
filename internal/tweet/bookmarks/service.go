package bookmarks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	bookmarksPathFormat      = "/2/users/%s/bookmarks"
	removeBookmarkPathFormat = "/2/users/%s/bookmarks/%s"
)

// Service coordinates Twitter bookmark operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the supplied client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("bookmarks: nil client")
	}
	return &Service{client: c}
}

// Add creates a bookmark for the specified tweet on behalf of the user.
func (s *Service) Add(ctx context.Context, userID, tweetID string) (BookmarkResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: nil service")
	}
	if userID == "" {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: user id is required")
	}
	if tweetID == "" {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: tweet id is required")
	}

	path := fmt.Sprintf(bookmarksPathFormat, userID)
	payload := BookmarkRequest{TweetID: tweetID}

	resp, err := s.client.Post(ctx, path, payload, nil)
	if err != nil {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return BookmarkResponse{}, rateLimits, err
	}

	var result BookmarkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return BookmarkResponse{}, rateLimits, fmt.Errorf("bookmarks: decode response: %w", err)
	}

	return result, rateLimits, nil
}

// Remove deletes a bookmark for the specified tweet on behalf of the user.
func (s *Service) Remove(ctx context.Context, userID, tweetID string) (BookmarkResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: nil service")
	}
	if userID == "" {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: user id is required")
	}
	if tweetID == "" {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: tweet id is required")
	}

	path := fmt.Sprintf(removeBookmarkPathFormat, userID, tweetID)

	resp, err := s.client.Delete(ctx, path, nil)
	if err != nil {
		return BookmarkResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return BookmarkResponse{}, rateLimits, err
	}

	var result BookmarkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return BookmarkResponse{}, rateLimits, fmt.Errorf("bookmarks: decode response: %w", err)
	}

	return result, rateLimits, nil
}

// List fetches the bookmarked tweets for the specified user.
func (s *Service) List(ctx context.Context, userID string, params map[string]string) (BookmarksListResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return BookmarksListResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: nil service")
	}
	if userID == "" {
		return BookmarksListResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("bookmarks: user id is required")
	}

	path := fmt.Sprintf(bookmarksPathFormat, userID)

	resp, err := s.client.Get(ctx, path, params)
	if err != nil {
		return BookmarksListResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return BookmarksListResponse{}, rateLimits, err
	}

	var result BookmarksListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return BookmarksListResponse{}, rateLimits, fmt.Errorf("bookmarks: decode response: %w", err)
	}

	return result, rateLimits, nil
}
