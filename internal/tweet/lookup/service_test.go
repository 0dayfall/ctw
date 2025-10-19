package lookup

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTweet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/tweets/tweet-1", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "300")
		w.Header().Set("x-rate-limit-remaining", "299")
		_, _ = w.Write([]byte(`{"data":[{"id":"tweet-1","text":"hello world"}]}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.GetTweet(context.Background(), "tweet-1", nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "tweet-1", resp.Data[0].ID)
	require.Equal(t, "hello world", resp.Data[0].Text)
	require.Equal(t, 300, rateLimits.Limit)
}

func TestGetTweets(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/tweets", r.URL.Path)
		require.Equal(t, "tweet-1,tweet-2", r.URL.Query().Get("ids"))

		w.Header().Set("x-rate-limit-limit", "300")
		w.Header().Set("x-rate-limit-remaining", "298")
		_, _ = w.Write([]byte(`{"data":[{"id":"tweet-1","text":"first"},{"id":"tweet-2","text":"second"}]}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.GetTweets(context.Background(), []string{"tweet-1", "tweet-2"}, nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 2)
	require.Equal(t, "tweet-1", resp.Data[0].ID)
	require.Equal(t, "tweet-2", resp.Data[1].ID)
	require.Equal(t, 300, rateLimits.Limit)
}
