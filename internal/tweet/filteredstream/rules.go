package tweet

import (
	"encoding/json"
	"log"

	common "github.com/0dayfall/ctw/internal/data"
	"github.com/0dayfall/ctw/internal/httphandler"
)

const (
	rules = "/2/tweets/search/stream"
)

var (
	rulesBaseUrl = common.APIurl + rules
)

func createRulesUrl(dryRun bool) (rulesUrl string) {
	rulesUrl = rulesBaseUrl + "/rules"
	if dryRun {
		rulesUrl += "?dry_run=true"
	}
	return
}

func AddRule(cmd AddCommand, dryRun bool) (jsonResponse RulesResponse, err error) {
	httpRequest := httphandler.CreatePostRequest(createRulesUrl(dryRun), cmd)
	httpResponse, err := httphandler.MakeRequest(httpRequest)
	defer httphandler.CloseBody(httpResponse.Body)
	httphandler.IsResponseOK(httpResponse)

	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return
}

func GetRules() (jsonResponse RulesResponse, err error) {
	httpRequest := httphandler.CreateGetRequest(createRulesUrl(false))
	httpResponse, err := httphandler.MakeRequest(httpRequest)
	defer httphandler.CloseBody(httpResponse.Body)

	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return
}
