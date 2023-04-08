package tweet

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0dayfall/ctw/internal/utils"
	"github.com/stretchr/testify/require"
)

func makeMap() map[string]string {
	m := make(map[string]string)
	m["tweet.fields"] = "created_at"
	m["expansions"] = "author_id"
	m["user.fields"] = "created_at"
	return m
}

func TestStreamsURL(t *testing.T) {
	require.Equal(t, "https://api.twitter.com/2/tweets/search/stream?"+
		"expansions=author_id&tweet.fields=created_at&user.fields=created_at",
		createStreamUrlWithFields(makeMap()))
}

func TestStream(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [
			  {
				"author_id": "2244994945",
				"created_at": "2022-09-14T19:00:55.000Z",
				"id": "1228393702244134912",
				"edit_history_tweet_ids": ["1228393702244134912"],
				"text": "What did the developer write in their Valentine’s card?\n  \nwhile(true) {\n    I = Love(You);  \n}"
			  },
			  {
				"author_id": "2244994945",
				"created_at": "2022-09-12T17:09:56.000Z",
				"id": "1227640996038684673",
				 "edit_history_tweet_ids": ["1227640996038684673"],
				"text": "Doctors: Googling stuff online does not make you a doctor\n\nDevelopers: https://t.co/mrju5ypPkb"
			  },
			  {
				"author_id": "2244994945",
				"created_at": "2022-09-27T20:26:41.000Z",
				"id": "1199786642791452673",
				"edit_history_tweet_ids": ["1199786642791452673"],
				"text": "C#"
			  }
			],
			"includes": {
			  "users": [
				{
				  "created_at": "2013-12-14T04:35:55.000Z",
				  "id": "2244994945",
				  "name": "Twitter Dev",
				  "username": "TwitterDev"
				}
			  ]
			}
		  }`))
		require.NoError(t, err)
	}))
	defer func() { testServer.Close() }()
	streamUrl = testServer.URL
	response, err := Stream(makeMap())
	require.NoError(t, err)
	require.True(t, utils.FindString(response, "Twitter Dev"))
}
