// Package timelines provides helpers for interacting with Twitter timeline endpoints.
package timelines

import "time"

// TimelineResponse captures paginated tweet timeline results.
type TimelineResponse struct {
	Data     []TweetData `json:"data"`
	Meta     Meta        `json:"meta"`
	Includes *Includes   `json:"includes,omitempty"`
}

// TweetData represents a tweet object in timeline responses.
type TweetData struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Meta provides pagination metadata for timeline responses.
type Meta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token,omitempty"`
	PreviousToken string `json:"previous_token,omitempty"`
	NewestID      string `json:"newest_id,omitempty"`
	OldestID      string `json:"oldest_id,omitempty"`
}

// Includes contains expanded objects referenced in the timeline.
type Includes struct {
	Users  []User      `json:"users,omitempty"`
	Tweets []TweetData `json:"tweets,omitempty"`
}

// User represents a user object in expanded includes.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
