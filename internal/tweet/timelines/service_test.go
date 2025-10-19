package timelines

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetUserTweets(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/users/123/tweets", r.URL.Path)
		require.Equal(t, "10", r.URL.Query().Get("max_results"))

		w.Header().Set("x-rate-limit-limit", "100")
		w.Header().Set("x-rate-limit-remaining", "99")
		w.Header().Set("x-rate-limit-reset", "1234567890")
		_, _ = w.Write([]byte(`{"data":[{"id":"1","text":"test tweet","author_id":"123"}],"meta":{"result_count":1,"newest_id":"1","oldest_id":"1"}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.GetUserTweets(context.Background(), "123", map[string]string{"max_results": "10"})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "1", resp.Data[0].ID)
	require.Equal(t, "test tweet", resp.Data[0].Text)
	require.Equal(t, 1, resp.Meta.ResultCount)
	require.Equal(t, 100, rateLimits.Limit)
	require.Equal(t, 99, rateLimits.Remaining)
}

func TestGetUserMentions(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/users/456/mentions", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "75")
		w.Header().Set("x-rate-limit-remaining", "74")
		_, _ = w.Write([]byte(`{"data":[{"id":"2","text":"@user mention"}],"meta":{"result_count":1}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.GetUserMentions(context.Background(), "456", nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "2", resp.Data[0].ID)
	require.Equal(t, "@user mention", resp.Data[0].Text)
	require.Equal(t, 75, rateLimits.Limit)
}

func TestGetReverseChronological(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/users/789/timelines/reverse_chronological", r.URL.Path)

		w.Header().Set("x-rate-limit-limit", "180")
		w.Header().Set("x-rate-limit-remaining", "179")
		response := map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "3", "text": "home timeline tweet"},
			},
			"meta": map[string]interface{}{
				"result_count": 1,
				"next_token":   "abc123",
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.GetReverseChronological(context.Background(), "789", nil)
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "3", resp.Data[0].ID)
	require.Equal(t, "home timeline tweet", resp.Data[0].Text)
	require.Equal(t, "abc123", resp.Meta.NextToken)
	require.Equal(t, 180, rateLimits.Limit)
	require.Equal(t, 179, rateLimits.Remaining)
}
