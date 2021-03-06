package tweet

import (
	"encoding/json"
	"log"

	common "github.com/0dayfall/ctw/data"
	"github.com/0dayfall/ctw/httphandler"
)

const (
	rules = "/2/tweets/search/stream"
)

func createStreamUrl() string {
	return common.APIurl + rules
}

func createRulesUrl(dryRun bool) string {
	if dryRun {
		return createStreamUrl() + "/rules?dry_run=true"
	}
	return createStreamUrl() + "/rules"
}

func GetRules() {
	url := createRulesUrl(false)
	httpRequest := httphandler.CreateGetRequest(url)
	response := httphandler.MakeRequest(httpRequest)
	defer response.Body.Close()

	var jsonResponse RulesResponse
	if err := json.NewDecoder(response.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	httphandler.PrettyPrint(jsonResponse)
}

func AddRule(cmd AddCommand, dryRun bool) {
	httpRequest := httphandler.CreatePostRequest(createRulesUrl(dryRun), cmd)
	httpResponse := httphandler.MakeRequest(httpRequest)
	defer httpResponse.Body.Close()
	httphandler.IsResponseOK(httpResponse)

	var jsonResponse RulesResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	httphandler.PrettyPrint(jsonResponse)
}

func Stream() {
	httpRequest := httphandler.CreateGetRequest(createStreamUrl())
	httpResponse := httphandler.MakeRequest(httpRequest)
	defer httpResponse.Body.Close()

	var jsonResponse RulesResponse
	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	httphandler.PrettyPrint(jsonResponse)
}
