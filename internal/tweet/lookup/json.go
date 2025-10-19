// Package lookup provides helpers for fetching tweets by ID.
package lookup

import "time"

// TweetLookupResponse captures single or multiple tweet lookups.
type TweetLookupResponse struct {
	Data     []Tweet   `json:"data"`
	Includes *Includes `json:"includes,omitempty"`
	Errors   []Error   `json:"errors,omitempty"`
}

// Tweet represents a tweet object.
type Tweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Includes contains expanded objects for the lookup.
type Includes struct {
	Users  []User  `json:"users,omitempty"`
	Tweets []Tweet `json:"tweets,omitempty"`
}

// User represents a user object in expanded includes.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

// Error captures API errors for individual tweet lookups.
type Error struct {
	Value        string `json:"value"`
	Detail       string `json:"detail"`
	Title        string `json:"title"`
	ResourceType string `json:"resource_type"`
	Parameter    string `json:"parameter"`
	Type         string `json:"type"`
}
