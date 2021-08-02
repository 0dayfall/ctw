package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/0dayfall/carboncopy/config"
)

func CreateGetRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+config.BearerToken)
	if config.UserAgent != "" {
		setUserAgent(req, config.UserAgent)
	}
	return req
}

func CreatePostRequest(url string, data interface{}) *http.Request {
	json, err := json.Marshal(&data)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(json))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.BearerToken)
	if config.UserAgent != "" {
		setUserAgent(req, config.UserAgent)
	}
	return req
}

func setUserAgent(req *http.Request, userAgent string) {
	req.Header.Set("User-Agent", userAgent)
}

type ErrorResponse struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
}

func IsResponseOK(response *http.Response) bool {
	responseOK := response.StatusCode > 199 && response.StatusCode < 300
	if !responseOK {
		log.Println(response.Status)
		var errorResponse ErrorResponse
		if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
			log.Println(err)
		}
		for _, error := range errorResponse.Errors {
			log.Println(error.Message)
		}
	}
	return responseOK
}

func IsRateLimitOK(resp *http.Response) (bool, int) {
	timeToReset, err := strconv.Atoi(resp.Header.Get("x-rate-limit-reset"))
	if err != nil {
		return false, -1
	}
	if resp.StatusCode == 429 {
		var jsonResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
			log.Println(err)
		}
		return false, timeToReset
	}
	return true, timeToReset
}

func MakeRequest(request *http.Request) *http.Response {
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func PrettyPrint(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		log.Println(string(b))
	}
}
