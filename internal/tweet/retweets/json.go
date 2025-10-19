// Package retweets provides helpers for interacting with Twitter retweet endpoints.
package retweets

import "time"

// RetweetRequest captures the payload required to retweet a tweet.
type RetweetRequest struct {
	TweetID string `json:"tweet_id"`
}

// RetweetResponse represents the ACK returned after creating or deleting a retweet.
type RetweetResponse struct {
	Data RetweetData `json:"data"`
}

// RetweetData provides the retweet status.
type RetweetData struct {
	Retweeted bool `json:"retweeted"`
}

// RetweetersResponse captures the list of users who retweeted a tweet.
type RetweetersResponse struct {
	Data []Retweeter `json:"data"`
	Meta Meta        `json:"meta"`
}

// Retweeter represents a user who retweeted.
type Retweeter struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Meta provides pagination metadata for retweeters.
type Meta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token,omitempty"`
	PreviousToken string `json:"previous_token,omitempty"`
}
