package tweet

import (
	"encoding/json"
	"log"
	"net/http"

	common "github.com/0dayfall/ctw/internal/data"
	httphandler "github.com/0dayfall/ctw/internal/httphandler"
	"github.com/0dayfall/ctw/internal/utils"
)

const (
	search = "/2/tweets/search/recent"
)

func createSearchTweetURL() string {
	return common.APIurl + search
}

type query struct {
	name  string
	value string
}

func addQuery(req *http.Request, queries []query) {
	q := req.URL.Query()
	for _, query := range queries {
		if len(query.value) > 0 {
			q.Add(query.name, query.value)
		}
	}
	req.URL.RawQuery = q.Encode()
	log.Println(req)
}

func SearchRecent(queryString string) (SearchRecentResponse, int, string) {
	url := createSearchTweetURL()
	req := httphandler.CreateGetRequest(url)
	addQuery(req, []query{
		{name: "query", value: queryString},
	})
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var jsonResponse SearchRecentResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	utils.PrettyPrint(jsonResponse)
	return jsonResponse, jsonResponse.Meta.ResultCount, jsonResponse.Meta.NextToken
}

func SearchRecentNextToken(queryString string, token string) (SearchRecentResponse, int, string) {
	url := createSearchTweetURL()
	req := httphandler.CreateGetRequest(url)
	addQuery(req, []query{
		{name: "query", value: queryString},
		{name: "pagination_token", value: token},
	})
	response := httphandler.MakeRequest(req)
	defer response.Body.Close()

	var jsonResponse SearchRecentResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	utils.PrettyPrint(jsonResponse)
	return jsonResponse, jsonResponse.Meta.ResultCount, jsonResponse.Meta.NextToken
}
