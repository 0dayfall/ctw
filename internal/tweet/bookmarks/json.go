// Package bookmarks provides helpers for interacting with Twitter bookmark endpoints.
package bookmarks

import "time"

// BookmarkRequest captures the payload required to bookmark a tweet.
type BookmarkRequest struct {
	TweetID string `json:"tweet_id"`
}

// BookmarkResponse represents the ACK returned after creating or removing a bookmark.
type BookmarkResponse struct {
	Data BookmarkData `json:"data"`
}

// BookmarkData provides the bookmark status.
type BookmarkData struct {
	Bookmarked bool `json:"bookmarked"`
}

// BookmarksListResponse captures the list of bookmarked tweets.
type BookmarksListResponse struct {
	Data []BookmarkedTweet `json:"data"`
	Meta Meta              `json:"meta"`
}

// BookmarkedTweet represents a single bookmarked tweet.
type BookmarkedTweet struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	AuthorID  string    `json:"author_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// Meta provides pagination metadata for bookmark listings.
type Meta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token,omitempty"`
	PreviousToken string `json:"previous_token,omitempty"`
}
