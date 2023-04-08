package tweet

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/0dayfall/ctw/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCountURL(t *testing.T) {
	require.True(t, "https://api.twitter.com/2/tweets/counts/recent?"+
		"granularity=day&query=from%3ATwitterDev" ==
		getRecentURL("from:TwitterDev", "day"))
}

func TestCountRecent(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [
				{
					"end": "2021-06-16T00:00:00.000Z",
					"start": "2021-06-15T00:00:00.000Z",
					"tweet_count": 0
				},
				{
					"end": "2021-06-17T00:00:00.000Z",
					"start": "2021-06-16T00:00:00.000Z",
					"tweet_count": 1
				},
				{
					"end": "2021-06-18T00:00:00.000Z",
					"start": "2021-06-17T00:00:00.000Z",
					"tweet_count": 2
				},
				{
					"end": "2021-06-19T00:00:00.000Z",
					"start": "2021-06-18T00:00:00.000Z",
					"tweet_count": 0
				},
				{
					"end": "2021-06-20T00:00:00.000Z",
					"start": "2021-06-19T00:00:00.000Z",
					"tweet_count": 0
				},
				{
					"end": "2021-06-21T00:00:00.000Z",
					"start": "2021-06-20T00:00:00.000Z",
					"tweet_count": 0
				},
				{
					"end": "2021-06-22T00:00:00.000Z",
					"start": "2021-06-21T00:00:00.000Z",
					"tweet_count": 1
				},
				{
					"end": "2021-06-23T00:00:00.000Z",
					"start": "2021-06-22T00:00:00.000Z",
					"tweet_count": 2
				}
			],
			"meta": {
				"total_tweet_count": 6
			}
		 }`))
		require.NoError(t, err)
	}))
	defer func() { testServer.Close() }()
	recentBaseURL = testServer.URL
	countResponse, err := GetRecentCount("from:TwitterDev", "day")
	require.NoError(t, err)
	require.EqualValues(t, 6, countResponse.Meta.TotalTweetCount)
	utils.PrettyPrint(countResponse)
}
