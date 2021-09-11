package config

import "os"

const (
	APIurl = "http://www.twitter.com"
)

var (
	BearerToken string = os.Getenv("BEARER_TOKEN")
	UserAgent   string
)

func Init(bearerToken string) {
	BearerToken = bearerToken
}

func SetUserAgent(userAgent string) {
	UserAgent = userAgent
}
