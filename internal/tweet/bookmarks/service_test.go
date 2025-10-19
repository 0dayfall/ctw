package bookmarks

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/2/users/123/bookmarks", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "49")
		_, _ = w.Write([]byte(`{"data":{"bookmarked":true}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.Add(context.Background(), "123", "tweet-1")
	require.NoError(t, err)
	require.True(t, resp.Data.Bookmarked)
	require.Equal(t, 50, rateLimits.Limit)
	require.Equal(t, 49, rateLimits.Remaining)
}

func TestRemove(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/2/users/456/bookmarks/tweet-2", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "48")
		_, _ = w.Write([]byte(`{"data":{"bookmarked":false}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.Remove(context.Background(), "456", "tweet-2")
	require.NoError(t, err)
	require.False(t, resp.Data.Bookmarked)
	require.Equal(t, 50, rateLimits.Limit)
}

func TestList(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/users/789/bookmarks", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "75")
		w.Header().Set("x-rate-limit-remaining", "74")
		_, _ = w.Write([]byte(`{"data":[{"id":"bm-1","text":"bookmarked tweet"}],"meta":{"result_count":1}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.List(context.Background(), "789", nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "bm-1", resp.Data[0].ID)
	require.Equal(t, "bookmarked tweet", resp.Data[0].Text)
	require.Equal(t, 75, rateLimits.Limit)
}
