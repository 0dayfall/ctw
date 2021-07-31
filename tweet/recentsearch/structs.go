package searchtweets

import "time"

type SearchRecentResponse struct {
	Data []Data `json:"data"`
	Meta Meta   `json:"meta"`
}

type ReferencedTweets struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Mentions struct {
	Start    int    `json:"start"`
	End      int    `json:"end"`
	Username string `json:"username"`
	ID       string `json:"id"`
}

type Annotations struct {
	Start          int     `json:"start"`
	End            int     `json:"end"`
	Probability    float64 `json:"probability"`
	Type           string  `json:"type"`
	NormalizedText string  `json:"normalized_text"`
}

type Entities struct {
	Mentions    []Mentions    `json:"mentions,omitempty"`
	Annotations []Annotations `json:"annotations,omitempty"`
	Urls        []Urls        `json:"urls,omitempty"`
	Hashtags    []Hashtags    `json:"hashtags,omitempty"`
}

type Hashtags struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Tag   string `json:"tag"`
}

type Images struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Urls struct {
	Start       int      `json:"start"`
	End         int      `json:"end"`
	URL         string   `json:"url"`
	ExpandedURL string   `json:"expanded_url"`
	DisplayURL  string   `json:"display_url"`
	Images      []Images `json:"images"`
	Status      int      `json:"status"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	UnwoundURL  string   `json:"unwound_url"`
}

type Data struct {
	ID                string             `json:"id"`
	ReferencedTweets  []ReferencedTweets `json:"referenced_tweets,omitempty"`
	Entities          Entities           `json:"entities,omitempty"`
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
