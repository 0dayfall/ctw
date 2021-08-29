package tweet

import (
	"encoding/json"
	"log"

	httphandler "github.com/0dayfall/ctw/httphandler"
)

func createSearchTweetURL() string {
	return "https://api.twitter.com/2/tweets/search/recent"
}

func SearchRecent(query string) (SearchRecentResponse, int, string) {
	url := createSearchTweetURL()
	req := httphandler.CreateGetRequest(url)
	q := req.URL.Query()
	q.Add("query", query)

	req.URL.RawQuery = q.Encode()
	log.Println(req)
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var jsonResponse SearchRecentResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	httphandler.PrettyPrint(jsonResponse)
	return jsonResponse, jsonResponse.Meta.ResultCount, jsonResponse.Meta.NextToken
}

func SearchRecentNextToken(query string, token string) (SearchRecentResponse, int, string) {
	url := createSearchTweetURL()
	req := httphandler.CreateGetRequest(url)
	q := req.URL.Query()
	q.Add("query", query)
	if token != "" {
		q.Add("pagination_token", token)
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
	return jsonResponse, jsonResponse.Meta.ResultCount, jsonResponse.Meta.NextToken
}
