package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/0dayfall/ctw/internal/client"
)

// Block prevents the target user from interacting with the source user.
func (s *Service) Block(ctx context.Context, sourceID, targetID string) (RelationshipResponse, client.RateLimitSnapshot, error) {
	path := fmt.Sprintf("/2/users/%s/blocking", sourceID)
	payload := map[string]string{"target_user_id": targetID}
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
		return RelationshipResponse{}, rateLimits, fmt.Errorf("userslookup: decode block response: %w", err)
	}

	return relationship, rateLimits, nil
}

// Unblock removes an existing block relationship.
func (s *Service) Unblock(ctx context.Context, sourceID, targetID string) (RelationshipResponse, client.RateLimitSnapshot, error) {
	path := fmt.Sprintf("/2/users/%s/blocking/%s", sourceID, targetID)
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
		return RelationshipResponse{}, rateLimits, fmt.Errorf("userslookup: decode unblock response: %w", err)
	}

	return relationship, rateLimits, nil
}

// Follow creates a follow relationship from sourceID to targetID.
func (s *Service) Follow(ctx context.Context, sourceID, targetID string) (RelationshipResponse, client.RateLimitSnapshot, error) {
	path := fmt.Sprintf("/2/users/%s/following", sourceID)
	payload := map[string]string{"target_user_id": targetID}
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
		return RelationshipResponse{}, rateLimits, fmt.Errorf("userslookup: decode follow response: %w", err)
	}

	return relationship, rateLimits, nil
}

// Unfollow removes a follow relationship.
func (s *Service) Unfollow(ctx context.Context, sourceID, targetID string) (RelationshipResponse, client.RateLimitSnapshot, error) {
	path := fmt.Sprintf("/2/users/%s/following/%s", sourceID, targetID)
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
		return RelationshipResponse{}, rateLimits, fmt.Errorf("userslookup: decode unfollow response: %w", err)
	}

	return relationship, rateLimits, nil
}

// RelationshipResponse captures the relationship mutation payloads.
type RelationshipResponse struct {
	Data RelationshipData `json:"data"`
}

// RelationshipData describes a relationship mutation result.
type RelationshipData struct {
	Blocking      bool `json:"blocking,omitempty"`
	Following     bool `json:"following,omitempty"`
	PendingFollow bool `json:"pending_follow,omitempty"`
}
