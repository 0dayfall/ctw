package BearerToken

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type bearer_token struct {
	token_type   string
	access_token string
}

func getBearerToken(api_key string, api_secret_key string) string {
	client := &http.Client{}
	req := getBearerTokenReq(api_key, api_secret_key)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var token bearer_token
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Fatal(err)
	}
	return token.access_token
}

func getBearerTokenReq(api_key string, api_secret_key string) *http.Request {
	req, err := http.NewRequest("POST", bearer_token_url, bytes.NewBuffer([]byte("grant_type=client_credentials")))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(api_key, api_secret_key)
	return req
}
