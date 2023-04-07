package config

import "os"

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
