package tweet

import (
	"encoding/json"
	"log"

	common "github.com/0dayfall/ctw/internal/data"
	httphandler "github.com/0dayfall/ctw/internal/httphandler"
)

const (
	search = "/2/tweets/search/recent"
)

var (
	recentSearchBaseURL = common.APIurl + search
)

func SearchRecent(queryString string) (searchRecentResponse SearchRecentResponse, resultCount int, nextToken string, err error) {
	req := httphandler.CreateGetRequest(recentSearchBaseURL)
	httphandler.AddQuery(req, map[string]string{"query": queryString})
	resp, err := httphandler.MakeRequest(req)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if err = json.NewDecoder(resp.Body).Decode(&searchRecentResponse); err != nil {
		log.Println(err)
	}
	resultCount = searchRecentResponse.Meta.ResultCount
	nextToken = searchRecentResponse.Meta.NextToken
	return
}

func SearchRecentNextToken(queryString string, token string) (searchRecentResponse SearchRecentResponse, resultCount int, nextToken string, err error) {
	req := httphandler.CreateGetRequest(recentSearchBaseURL)
	httphandler.AddQuery(req, map[string]string{
		"query":            queryString,
		"pagination_token": token,
	})
	httpResponse, err := httphandler.MakeRequest(req)
	defer httphandler.CloseBody(httpResponse.Body)

	if err = json.NewDecoder(httpResponse.Body).Decode(&searchRecentResponse); err != nil {
		log.Println(err)
	}
	resultCount = searchRecentResponse.Meta.ResultCount
	nextToken = searchRecentResponse.Meta.NextToken
	return
}
