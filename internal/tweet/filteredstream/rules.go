package tweet

import (
	"encoding/json"
	"log"
	"net/url"

	common "github.com/0dayfall/ctw/internal/data"
	"github.com/0dayfall/ctw/internal/httphandler"
)

const (
	rules = "/2/tweets/search/stream"
)

var (
	rulesUrl = common.APIurl + rules
)

func createStreamUrl() string {
	return rulesUrl
}

func createStreamUrlWithFields(fields map[string]string) string {
	params := url.Values{}
	for key, value := range fields {
		params.Add(key, value)
	}
	return rulesUrl + "?" + params.Encode()
}

func createRulesUrl(dryRun bool) (rulesUrl string) {
	rulesUrl = createStreamUrl() + "/rules"
	if dryRun {
		rulesUrl += "?dry_run=true"
	}
	return
}

func GetRules() (jsonResponse RulesResponse, err error) {
	httpRequest := httphandler.CreateGetRequest(createRulesUrl(false))
	response := httphandler.MakeRequest(httpRequest)
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return
}

func AddRule(cmd AddCommand, dryRun bool) (jsonResponse RulesResponse, err error) {
	httpRequest := httphandler.CreatePostRequest(createRulesUrl(dryRun), cmd)
	httpResponse := httphandler.MakeRequest(httpRequest)
	defer httpResponse.Body.Close()
	httphandler.IsResponseOK(httpResponse)

	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return
}

func Stream(fields map[string]string) (jsonResponse interface{}, err error) {
	httpRequest := httphandler.CreateGetRequest(createStreamUrlWithFields(fields))
	httpResponse := httphandler.MakeRequest(httpRequest)
	defer httpResponse.Body.Close()

	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return
}
