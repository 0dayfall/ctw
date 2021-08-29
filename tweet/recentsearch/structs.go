package tweet

import (
	"time"

	common "github.com/0dayfall/ctw/data"
)

type SearchRecentResponse struct {
	Data []Data `json:"data"`
	Meta Meta   `json:"meta"`
}

type ReferencedTweets struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Data struct {
	ID                string             `json:"id"`
	ReferencedTweets  []ReferencedTweets `json:"referenced_tweets,omitempty"`
	Entities          common.Entities    `json:"entities,omitempty"`
	CreatedAt         time.Time          `json:"created_at"`
	PossiblySensitive bool               `json:"possibly_sensitive"`
	Text              string             `json:"text"`
	Source            string             `json:"source"`
	Lang              string             `json:"lang"`
	AuthorID          string             `json:"author_id"`
	InReplyToUserID   string             `json:"in_reply_to_user_id,omitempty"`
	ConversationID    string             `json:"conversation_id"`
}

type Meta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}
