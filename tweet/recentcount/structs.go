package tweet

import "time"

type CountResponse struct {
	Data []Data `json:"data"`
	Meta Meta   `json:"meta"`
}
type Data struct {
	End        time.Time `json:"end"`
	Start      time.Time `json:"start"`
	TweetCount int       `json:"tweet_count"`
}
type Meta struct {
	TotalTweetCount int `json:"total_tweet_count"`
}
