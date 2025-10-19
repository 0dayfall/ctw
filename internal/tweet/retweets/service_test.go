package retweets

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRetweet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/2/users/123/retweets", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "49")
		_, _ = w.Write([]byte(`{"data":{"retweeted":true}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.Retweet(context.Background(), "123", "tweet-1")
	require.NoError(t, err)
	require.True(t, resp.Data.Retweeted)
	require.Equal(t, 50, rateLimits.Limit)
	require.Equal(t, 49, rateLimits.Remaining)
}

func TestUnretweet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/2/users/456/retweets/tweet-2", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "48")
		_, _ = w.Write([]byte(`{"data":{"retweeted":false}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.Unretweet(context.Background(), "456", "tweet-2")
	require.NoError(t, err)
	require.False(t, resp.Data.Retweeted)
	require.Equal(t, 50, rateLimits.Limit)
}

func TestListRetweeters(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/tweets/tweet-3/retweeted_by", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "75")
		w.Header().Set("x-rate-limit-remaining", "74")
		_, _ = w.Write([]byte(`{"data":[{"id":"u1","name":"User One","username":"user1"}],"meta":{"result_count":1}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.ListRetweeters(context.Background(), "tweet-3", nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "u1", resp.Data[0].ID)
	require.Equal(t, "User One", resp.Data[0].Name)
	require.Equal(t, 75, rateLimits.Limit)
}
