package httphandler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/0dayfall/ctw/internal/config"
	"github.com/0dayfall/ctw/internal/utils"
)

func CreateGetRequest(url string) *http.Request {
	httpRequest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	setupHeaders(httpRequest)
	return httpRequest
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

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func IsResponseOK(httpResponse *http.Response) bool {
	responseOK := httpResponse.StatusCode > 199 && httpResponse.StatusCode < 300
	if !responseOK {
		log.Println(httpResponse.Status)
		var errorResponse ErrorResponse
		if err := json.NewDecoder(httpResponse.Body).Decode(&errorResponse); err != nil {
			log.Fatal(err)
		}
		for _, error := range errorResponse.Errors {
			utils.PrettyPrint(error)
		}
	}
	return responseOK
}

func IsRateLimitResetOK(httpResponse *http.Response) (bool, int) {
	timeToReset, err := strconv.Atoi(httpResponse.Header.Get("x-rate-limit-reset"))
	if err != nil {
		log.Println(err)
		return false, -1
	}
	if httpResponse.StatusCode == http.StatusTooManyRequests {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(httpResponse.Body).Decode(&errorResponse); err != nil {
			log.Fatal(err)
		}
		for _, v := range errorResponse.Errors {
			log.Printf("The HTTP status was %d %s, error code: %d, message: %s\n", http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests), v.Code, v.Message)
		}
		return false, timeToReset
	}
	return true, timeToReset
}

func GetRateLimitLimit(httpResponse *http.Response) (rateLimitLimit int) {
	rateLimitLimit, err := strconv.Atoi(httpResponse.Header.Get("x-rate-limit-limit"))
	if err != nil {
		return -1
	}
	return
}

func GetRateLimitRemaining(httpResponse *http.Response) (rateLimitRemaining int) {
	rateLimitRemaining, err := strconv.Atoi(httpResponse.Header.Get("x-rate-limit-remaining"))
	if err != nil {
		return -1
	}
	return
}

func MakeRequest(httpRequest *http.Request) (httpResponse *http.Response, err error) {
	client := http.Client{
		Timeout: 60 * time.Second,
	}
	httpResponse, err = client.Do(httpRequest)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func setupHeaders(httpRequest *http.Request) {
	addContentType(httpRequest)
	addBearerToken(httpRequest)
	setUserAgent(httpRequest)
	addAcceptEncoding(httpRequest)
}

func addContentType(httpRequest *http.Request) {
	httpRequest.Header.Add("Content-type", "application/json")
}

func addBearerToken(httpRequest *http.Request) {
	httpRequest.Header.Add("Authorization", "Bearer "+config.BearerToken)
}

func setUserAgent(httpRequest *http.Request) {
	if config.UserAgent != "" {
		httpRequest.Header.Set("User-Agent", config.UserAgent)
	}
}

func addAcceptEncoding(httpRequest *http.Request) {
	httpRequest.Header.Add("Accept-Encoding", "gzip")
}

func AddQuery(httpRequest *http.Request, queries map[string]string) {
	q := httpRequest.URL.Query()
	for k, v := range queries {
		q.Add(k, v)
	}
	httpRequest.URL.RawQuery = q.Encode()
	log.Println(httpRequest)
}

func CloseBody(closer io.ReadCloser) {
	err := closer.Close()
	if err != nil {
		log.Println(err)
	}
}
