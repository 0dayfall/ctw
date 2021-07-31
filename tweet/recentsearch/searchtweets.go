package searchtweets

import (
	"encoding/json"
	"log"

	httphandler "github.com/0dayfall/carboncopy/httphandler"
)

func createSearchTweetURL() string {
	return "https://api.twitter.com/2/tweets/search/recent"
}

func SearchRecent(query string) {
	var nextToken = ""
	for {
		url := createSearchTweetURL()
		req := httphandler.CreateGetRequest(url)
		q := req.URL.Query()
		q.Add("query", query)
		if nextToken != "" {
			q.Add("pagination_token", nextToken)
		}
		req.URL.RawQuery = q.Encode()
		log.Println(req)
		response := httphandler.MakeRequest(req)
		defer response.Body.Close()

		var jsonResponse SearchRecentResponse
		if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
			log.Println(err)
		}
		httphandler.PrettyPrint(jsonResponse)
		if jsonResponse.Meta.NextToken == "" {
			break
		}
		nextToken = jsonResponse.Meta.NextToken
	}
}
