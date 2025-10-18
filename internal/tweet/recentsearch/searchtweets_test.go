package tweet

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchRecent(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/2/tweets/search/recent", req.URL.Path)
		require.Equal(t, "ericsson lang:sv", req.URL.Query().Get("query"))
		res.Header().Set("x-rate-limit-limit", "450")
		res.Header().Set("x-rate-limit-remaining", "440")
		res.Header().Set("x-rate-limit-reset", "120")
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [{
				"id": "1",
				"text": "tweet",
				"author_id": "42",
				"created_at": "2024-01-01T00:00:00.000Z",
				"conversation_id": "99",
				"possibly_sensitive": false,
				"source": "api",
				"lang": "sv"
			}],
			"meta": {
				"result_count": 1,
				"next_token": "abc"
			}
		}`))
		require.NoError(t, err)
	})

	response, limits, err := service.SearchRecent(context.Background(), "ericsson lang:sv", nil)
	require.NoError(t, err)
	require.Len(t, response.Data, 1)
	require.Equal(t, 1, response.Meta.ResultCount)
	require.Equal(t, "abc", response.Meta.NextToken)
	require.Equal(t, 450, limits.Limit)
	require.Equal(t, 440, limits.Remaining)
	require.Equal(t, 120, limits.Reset)
}

func TestSearchRecentNextToken(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "abc", req.URL.Query().Get("pagination_token"))
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [],
			"meta": {
				"result_count": 0
			}
		}`))
		require.NoError(t, err)
	})

	response, _, err := service.SearchRecentNextToken(context.Background(), "ericsson lang:sv", "abc")
	require.NoError(t, err)
	require.Zero(t, response.Meta.ResultCount)
}
