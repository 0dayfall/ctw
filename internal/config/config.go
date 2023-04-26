package config

import "os"

var (
	BearerToken string = os.Getenv("BEARER_TOKEN")
	UserAgent   string = "CERN-LineMode/2.15 libwww/2.17b3"
)

func Init(bearerToken string) {
	BearerToken = bearerToken
}

func SetUserAgent(userAgent string) {
	UserAgent = userAgent
}
