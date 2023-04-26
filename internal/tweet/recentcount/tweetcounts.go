package tweet

import (
	"encoding/json"
	"log"
	"net/url"

	common "github.com/0dayfall/ctw/internal/data"
	"github.com/0dayfall/ctw/internal/httphandler"
)

const (
	recent = "/2/tweets/counts/recent"
	all    = "/2/tweets/counts/all"
)

var (
	recentBaseURL   = common.APIurl + recent
	allTweetBaseURL = common.APIurl + all
)

func getRecentURL(query string, granularity string) string {
	params := url.Values{}
	params.Add("query", query)
	params.Add("granularity", granularity)
	return recentBaseURL + "?" + params.Encode()
}

func GetRecentCount(query string, granularity string) (countResponse CountResponse, err error) {
	req := httphandler.CreateGetRequest(getRecentURL(query, granularity))

	httpResponse, err := httphandler.MakeRequest(req)
	defer httphandler.CloseBody(httpResponse.Body)

	if !httphandler.IsResponseOK(httpResponse) {
		return
	}

	if err = json.NewDecoder(httpResponse.Body).Decode(&countResponse); err != nil {
		log.Println(err)
	}

	return
}
