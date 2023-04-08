package httphandler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/0dayfall/ctw/internal/config"
	"github.com/0dayfall/ctw/internal/utils"
)

func CreateGetRequest(url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	setupHeaders(req)
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
	setupHeaders(req)
	return req
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
			utils.PrettyPrint(error)
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

func MakeRequest(request *http.Request) (response *http.Response, err error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	response, err = client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func setupHeaders(req *http.Request) {
	addContentType(req)
	addBearerToken(req)
	setUserAgent(req)
}

func addContentType(req *http.Request) {
	req.Header.Add("Content-type", "application/json")
}

func addBearerToken(req *http.Request) {
	req.Header.Add("Authorization", "Bearer "+config.BearerToken)
}

func setUserAgent(req *http.Request) {
	if config.UserAgent != "" {
		req.Header.Set("User-Agent", config.UserAgent)
	}
}

func AddQuery(req *http.Request, queries map[string]string) {
	q := req.URL.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	log.Println(req)
}
