package publish

// CreateTweetRequest represents the minimal payload required to create a tweet.
type CreateTweetRequest struct {
	Text string `json:"text,omitempty"`
}

// CreateTweetResponse captures the response payload for POST /2/tweets.
type CreateTweetResponse struct {
	Data TweetData `json:"data"`
}

// TweetData contains the created tweet metadata returned by the Twitter API.
type TweetData struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// DeleteTweetResponse captures the response payload from DELETE /2/tweets/:id.
type DeleteTweetResponse struct {
	Data DeleteData `json:"data"`
}

// DeleteData indicates whether the delete operation succeeded.
type DeleteData struct {
	Deleted bool `json:"deleted"`
}
