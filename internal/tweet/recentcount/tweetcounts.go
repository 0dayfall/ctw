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

	resp, err := httphandler.MakeRequest(req)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	if !httphandler.IsResponseOK(resp) {
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&countResponse); err != nil {
		log.Println(err)
	}

	return
}
