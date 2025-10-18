// Package user provides user lookup and relationship management helpers for
// the Twitter API v2.
package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/client"
	common "github.com/0dayfall/ctw/internal/data"
)

const (
	usersPath          = "/2/users"
	userByIDPath       = "/2/users/%s"
	usersByUsername    = "/2/users/by"
	userByUsernamePath = "/2/users/by/username/%s"
)

// Service coordinates user lookup operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the provided client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("userslookup: nil client")
	}
	return &Service{client: c}
}

// LookupID fetches a single user by ID.
func (s *Service) LookupID(ctx context.Context, id string, params map[string]string) (User, client.RateLimitSnapshot, error) {
	path := fmt.Sprintf(userByIDPath, id)
	return s.fetchSingle(ctx, path, params)
}

// LookupUsername fetches a single user by username.
func (s *Service) LookupUsername(ctx context.Context, username string, params map[string]string) (User, client.RateLimitSnapshot, error) {
	path := fmt.Sprintf(userByUsernamePath, username)
	return s.fetchSingle(ctx, path, params)
}

// LookupIDs fetches multiple users by comma-separated IDs.
func (s *Service) LookupIDs(ctx context.Context, ids []string, params map[string]string) ([]User, client.RateLimitSnapshot, error) {
	qp := map[string]string{
		"ids": strings.Join(ids, ","),
	}
	for k, v := range params {
		if k == "ids" {
			continue
		}
		qp[k] = v
	}

	resp, err := s.client.Get(ctx, usersPath, qp)
	if err != nil {
		return nil, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return nil, rateLimits, err
	}

	var payload UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, rateLimits, fmt.Errorf("userslookup: decode ids response: %w", err)
	}

	return payload.Data, rateLimits, nil
}

// LookupUsernames fetches multiple users by username.
func (s *Service) LookupUsernames(ctx context.Context, usernames []string, params map[string]string) ([]User, client.RateLimitSnapshot, error) {
	qp := map[string]string{
		"usernames": strings.Join(usernames, ","),
	}
	for k, v := range params {
		if k == "usernames" {
			continue
		}
		qp[k] = v
	}

	resp, err := s.client.Get(ctx, usersByUsername, qp)
	if err != nil {
		return nil, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return nil, rateLimits, err
	}

	var payload UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, rateLimits, fmt.Errorf("userslookup: decode usernames response: %w", err)
	}

	return payload.Data, rateLimits, nil
}

func (s *Service) fetchSingle(ctx context.Context, path string, params map[string]string) (User, client.RateLimitSnapshot, error) {
	resp, err := s.client.Get(ctx, path, params)
	if err != nil {
		return User{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return User{}, rateLimits, err
	}

	var payload UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return User{}, rateLimits, fmt.Errorf("userslookup: decode response: %w", err)
	}

	return payload.Data, rateLimits, nil
}

// User describes the Twitter user payload returned by the API.
type User struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	UserName        string          `json:"username"`
	CreatedAt       string          `json:"created_at"`
	Description     string          `json:"description"`
	Entities        common.Entities `json:"entities"`
	Location        string          `json:"location"`
	PinnedTweetID   string          `json:"pinned_tweet_id"`
	ProfileImageURL string          `json:"profile_image_url"`
	Protected       bool            `json:"protected"`
	PublicMetrics   UserMetrics     `json:"public_metrics"`
	URL             string          `json:"url"`
	Verified        bool            `json:"verified"`
	WithHeld        WithHeld        `json:"withheld"`
}

// WithHeld captures withheld information for a user.
type WithHeld struct {
	Copyright    bool     `json:"copyright"`
	CountryCodes []string `json:"country_codes"`
}

// UserMetrics contains activity metrics.
type UserMetrics struct {
	Followers int `json:"followers_count"`
	Following int `json:"following_count"`
	Tweets    int `json:"tweet_count"`
	Listed    int `json:"listed_count"`
}

// UserResponse wraps a single user payload.
type UserResponse struct {
	Data User `json:"data"`
}

// UsersResponse wraps multiple user payloads.
type UsersResponse struct {
	Data []User `json:"data"`
}
