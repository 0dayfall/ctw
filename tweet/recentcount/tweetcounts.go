package tweetcounts

import (
	"encoding/json"
	"log"
	"os"

	"github.com/0dayfall/carboncopy/httphandler"
)

var (
	bearerToken = os.Getenv("BEARER_TOKEN")
)

func createRecentTweetCountsUrl() string {
	return "https://api.twitter.com/2/tweets/counts/recent"
}

func createAllTweetCountsUrl() string {
	return "https://api.twitter.com/2/tweets/counts/all"
}

func GetRecentCount(query string, granularity string) CountResponse {
	url := createRecentTweetCountsUrl()
	req := httphandler.CreateGetRequest(url)
	q := req.URL.Query()
	q.Add("query", query)
	q.Add("granularity", granularity)
	req.URL.RawQuery = q.Encode()
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()
	if !httphandler.IsResponseOK(response) {
		return CountResponse{}
	}
	var jsonResponse CountResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return jsonResponse
}
