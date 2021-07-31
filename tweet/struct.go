package carboncopy

import "time"

type TweetObject struct {
	Data Data `json:"data"`
}

type Data struct {
	AuthorID          string    `json:"author_id"`
	CreatedAt         time.Time `json:"created_at"`
	Entities          Entities  `json:"entities"`
	ID                string    `json:"id"`
	Lang              string    `json:"lang"`
	PossiblySensitive bool      `json:"possibly_sensitive"`
	Source            string    `json:"source"`
	Text              string    `json:"text"`
}

type Entities struct {
	Urls []Urls `json:"urls"`
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

type Images struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
