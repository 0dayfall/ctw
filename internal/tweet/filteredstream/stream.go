package tweet

import (
	"encoding/json"
	"log"
	"net/url"

	common "github.com/0dayfall/ctw/internal/data"
	"github.com/0dayfall/ctw/internal/httphandler"
)

var (
	streamUrl = common.APIurl + rules
)

func createStreamUrlWithFields(fields map[string]string) string {
	params := url.Values{}
	for key, value := range fields {
		params.Add(key, value)
	}
	return streamUrl + "?" + params.Encode()
}

func Stream(fields map[string]string) (jsonResponse map[string]interface{}, err error) {
	httpRequest := httphandler.CreateGetRequest(createStreamUrlWithFields(fields))
	httpResponse, err := httphandler.MakeRequest(httpRequest)
	defer httphandler.CloseBody(httpResponse.Body)

	if err := json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		log.Println(err)
	}
	return
}
