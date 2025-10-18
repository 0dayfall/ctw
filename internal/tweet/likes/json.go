// Package likes provides helpers for interacting with Tweet like endpoints.
package likes

import "time"

// RelationshipResponse captures like/unlike mutation responses.
type RelationshipResponse struct {
	Data RelationshipData `json:"data"`
}

// RelationshipData describes the like status returned by the API.
type RelationshipData struct {
	Liked bool `json:"liked"`
}

// LikedTweetsResponse captures GET /2/users/:id/liked_tweets payloads.
type LikedTweetsResponse struct {
	Data []LikedTweet `json:"data"`
	Meta Meta         `json:"meta"`
}

// LikedTweet represents a single liked tweet record.
type LikedTweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Meta provides pagination metadata for liked tweets.
type Meta struct {
	NextToken   string `json:"next_token"`
	ResultCount int    `json:"result_count"`
}
