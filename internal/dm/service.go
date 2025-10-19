package dm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/0dayfall/ctw/internal/client"
)

const (
	dmEventsPath             = "/2/dm_events"
	dmConversationPathFormat = "/2/dm_conversations/%s/messages"
	dmWithUserPathFormat     = "/2/dm_conversations/with/%s/messages"
)

// Service coordinates Twitter direct message operations.
type Service struct {
	client *client.Client
}

// NewService constructs a Service backed by the supplied client.
func NewService(c *client.Client) *Service {
	if c == nil {
		panic("dm: nil client")
	}
	return &Service{client: c}
}

// SendToUser starts or continues a 1:1 conversation with the given participant.
func (s *Service) SendToUser(ctx context.Context, participantID string, req SendDMRequest) (SendDMResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return SendDMResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: nil service")
	}
	if participantID == "" {
		return SendDMResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: participant id is required")
	}
	if strings.TrimSpace(req.Text) == "" {
		return SendDMResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: text is required")
	}

	path := fmt.Sprintf(dmWithUserPathFormat, participantID)
	return s.send(ctx, path, req)
}

// SendToConversation posts a message to an existing DM conversation.
func (s *Service) SendToConversation(ctx context.Context, conversationID string, req SendDMRequest) (SendDMResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return SendDMResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: nil service")
	}
	if conversationID == "" {
		return SendDMResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: conversation id is required")
	}
	if strings.TrimSpace(req.Text) == "" {
		return SendDMResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: text is required")
	}

	path := fmt.Sprintf(dmConversationPathFormat, conversationID)
	return s.send(ctx, path, req)
}

func (s *Service) send(ctx context.Context, path string, req SendDMRequest) (SendDMResponse, client.RateLimitSnapshot, error) {
	resp, err := s.client.Post(ctx, path, req, nil)
	if err != nil {
		return SendDMResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return SendDMResponse{}, rateLimits, err
	}

	var payload SendDMResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return SendDMResponse{}, rateLimits, fmt.Errorf("dm: decode response: %w", err)
	}

	return payload, rateLimits, nil
}

// ListEvents returns DM events for the authenticated user using optional query params.
func (s *Service) ListEvents(ctx context.Context, params map[string]string) (DMEventsResponse, client.RateLimitSnapshot, error) {
	if s == nil {
		return DMEventsResponse{}, client.RateLimitSnapshot{}, fmt.Errorf("dm: nil service")
	}

	resp, err := s.client.Get(ctx, dmEventsPath, params)
	if err != nil {
		return DMEventsResponse{}, client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return DMEventsResponse{}, rateLimits, err
	}

	var payload DMEventsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return DMEventsResponse{}, rateLimits, fmt.Errorf("dm: decode response: %w", err)
	}

	return payload, rateLimits, nil
}

// DeleteEvent removes a DM event by identifier, returning any rate limit headers.
func (s *Service) DeleteEvent(ctx context.Context, eventID string) (client.RateLimitSnapshot, error) {
	if s == nil {
		return client.RateLimitSnapshot{}, fmt.Errorf("dm: nil service")
	}
	if eventID == "" {
		return client.RateLimitSnapshot{}, fmt.Errorf("dm: event id is required")
	}

	path := fmt.Sprintf("%s/%s", dmEventsPath, eventID)
	resp, err := s.client.Delete(ctx, path, nil)
	if err != nil {
		return client.RateLimitSnapshot{}, err
	}
	defer client.SafeClose(resp.Body)

	rateLimits := client.ParseRateLimits(resp)
	if err := client.CheckResponse(resp); err != nil {
		return rateLimits, err
	}

	return rateLimits, nil
}
