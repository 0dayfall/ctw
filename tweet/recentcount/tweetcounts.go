package tweet

import (
	"encoding/json"
	"log"

	common "github.com/0dayfall/ctw/data"
	"github.com/0dayfall/ctw/httphandler"
)

const (
	recentURL   = "/2/tweets/counts/recent"
	countAllURL = "/2/tweets/counts/all"
)

func createRecentTweetCountsUrl() string {
	return common.APIurl + recentURL
}

func createAllTweetCountsUrl() string {
	return common.APIurl + countAllURL
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
	var countResponse CountResponse
	if err := json.NewDecoder(response.Body).Decode(&countResponse); err != nil {
		log.Println(err)
	}
	return countResponse
}
