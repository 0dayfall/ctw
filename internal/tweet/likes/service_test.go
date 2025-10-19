package likes

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLikeTweet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/2/users/123/likes", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "49")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"liked":true}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.LikeTweet(context.Background(), "123", "tweet-1")
	require.NoError(t, err)
	require.True(t, resp.Data.Liked)
	require.Equal(t, 50, rateLimits.Limit)
	require.Equal(t, 49, rateLimits.Remaining)
}

func TestUnlikeTweet(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/2/users/456/likes/tweet-2", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "48")
		_, _ = w.Write([]byte(`{"data":{"liked":false}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.UnlikeTweet(context.Background(), "456", "tweet-2")
	require.NoError(t, err)
	require.False(t, resp.Data.Liked)
	require.Equal(t, 50, rateLimits.Limit)
	require.Equal(t, 48, rateLimits.Remaining)
}

func TestListLikedTweets(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/users/789/liked_tweets", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "75")
		w.Header().Set("x-rate-limit-remaining", "74")
		_, _ = w.Write([]byte(`{"data":[{"id":"like-1","text":"liked tweet"}],"meta":{"result_count":1,"next_token":"token-abc"}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.ListLikedTweets(context.Background(), "789", nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "like-1", resp.Data[0].ID)
	require.Equal(t, "liked tweet", resp.Data[0].Text)
	require.Equal(t, "token-abc", resp.Meta.NextToken)
	require.Equal(t, 75, rateLimits.Limit)
}
