package dm

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendToUser(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/2/dm_conversations/with/123/messages", r.URL.Path)

		var payload SendDMRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "hello dm", payload.Text)

		w.Header().Set("x-rate-limit-limit", "50")
		w.Header().Set("x-rate-limit-remaining", "49")
		w.Header().Set("x-rate-limit-reset", "100")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"data":{"dm_event_id":"event-1"}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.SendToUser(context.Background(), "123", SendDMRequest{Text: "hello dm"})
	require.NoError(t, err)
	require.Equal(t, "event-1", resp.Data.DMEventID)
	require.Equal(t, 50, rateLimits.Limit)
	require.Equal(t, 49, rateLimits.Remaining)
	require.Equal(t, 100, rateLimits.Reset)
}

func TestListEvents(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/2/dm_events", r.URL.Path)
		require.Equal(t, "token-1", r.URL.Query().Get("pagination_token"))

		w.Header().Set("x-rate-limit-limit", "10")
		w.Header().Set("x-rate-limit-remaining", "9")
		w.Header().Set("x-rate-limit-reset", "200")
		_, _ = w.Write([]byte(`{"data":[{"id":"event-2","text":"hey there","event_type":"message_create","conversation_id":"conv-1","dm_conversation_id":"123-456","sender_id":"123","created_at":"2023-01-01T00:00:00Z"}],"meta":{"result_count":1,"next_token":"token-2"}}`))
	}

	service := newTestService(t, handler)

	resp, rateLimits, err := service.ListEvents(context.Background(), map[string]string{"pagination_token": "token-1"})
	require.NoError(t, err)
	require.Len(t, resp.Data, 1)
	require.Equal(t, "event-2", resp.Data[0].ID)
	require.Equal(t, 1, resp.Meta.ResultCount)
	require.Equal(t, "token-2", resp.Meta.NextToken)
	require.Equal(t, 10, rateLimits.Limit)
	require.Equal(t, 9, rateLimits.Remaining)
}

func TestDeleteEvent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/2/dm_events/event-99", r.URL.Path)
		w.Header().Set("x-rate-limit-limit", "5")
		w.Header().Set("x-rate-limit-remaining", "4")
		w.Header().Set("x-rate-limit-reset", "300")
		w.WriteHeader(http.StatusNoContent)
	}

	service := newTestService(t, handler)

	rateLimits, err := service.DeleteEvent(context.Background(), "event-99")
	require.NoError(t, err)
	require.Equal(t, 5, rateLimits.Limit)
	require.Equal(t, 4, rateLimits.Remaining)
	require.Equal(t, 300, rateLimits.Reset)
}
