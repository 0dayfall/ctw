// Package dm provides helpers for interacting with Twitter direct message endpoints.
package dm

import "time"

// SendDMRequest captures the payload required to create a direct message.
type SendDMRequest struct {
	Text string `json:"text"`
}

// SendDMResponse represents the ACK returned after creating a DM event.
type SendDMResponse struct {
	Data SendDMData `json:"data"`
}

// SendDMData provides identifiers for the newly created DM event.
type SendDMData struct {
	DMEventID string `json:"dm_event_id"`
}

// DMEventsResponse captures the timeline-style listing of DM events.
type DMEventsResponse struct {
	Data []DMEvent `json:"data"`
	Meta DMMeta    `json:"meta"`
}

// DMEvent represents a single DM event returned by Twitter.
type DMEvent struct {
	ID               string    `json:"id"`
	Text             string    `json:"text"`
	EventType        string    `json:"event_type"`
	ConversationID   string    `json:"conversation_id"`
	DMConversationID string    `json:"dm_conversation_id"`
	SenderID         string    `json:"sender_id"`
	CreatedAt        time.Time `json:"created_at"`
}

// DMMeta contains pagination metadata for DM listings.
type DMMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}
