package tweet

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRecentCount(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/2/tweets/counts/recent", req.URL.Path)
		require.Equal(t, "from:TwitterDev", req.URL.Query().Get("query"))
		require.Equal(t, "day", req.URL.Query().Get("granularity"))
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [
				{
					"end": "2021-06-16T00:00:00.000Z",
					"start": "2021-06-15T00:00:00.000Z",
					"tweet_count": 0
				}
			],
			"meta": {
				"total_tweet_count": 0
			}
		}`))
		require.NoError(t, err)
	})

	response, _, err := service.GetRecentCount(context.Background(), "from:TwitterDev", "day", nil)
	require.NoError(t, err)
	require.Zero(t, response.Meta.TotalTweetCount)
}
