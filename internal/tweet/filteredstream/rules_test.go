package tweet

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddRuleDryRun(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/2/tweets/search/stream/rules", req.URL.Path)
		require.Equal(t, "true", req.URL.Query().Get("dry_run"))
		res.Header().Set("x-rate-limit-limit", "450")
		res.Header().Set("x-rate-limit-remaining", "449")
		res.Header().Set("x-rate-limit-reset", "60")
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [{
				"value": "cat has:images",
				"tag": "cats with images",
				"id": "1273026480692322304"
			}],
			"meta": {
				"sent": "2020-06-16T22:55:39.356Z",
				"summary": {
					"created": 1,
					"not_created": 0,
					"valid": 1,
					"invalid": 0
				}
			}
		}`))
		require.NoError(t, err)
	})

	response, rateLimits, err := service.AddRule(context.Background(), AddCommand{
		Add: []Add{{
			Value: "cat has:images",
			Tag:   "cats with images",
		}},
	}, true)

	require.NoError(t, err)
	require.EqualValues(t, 1, response.Meta.Summary.Created)
	require.Equal(t, 450, rateLimits.Limit)
	require.Equal(t, 449, rateLimits.Remaining)
	require.Equal(t, 60, rateLimits.Reset)
}

func TestAddRule(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/2/tweets/search/stream/rules", req.URL.Path)
		require.Empty(t, req.URL.Query().Get("dry_run"))
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [{
				"value": "cat has:images",
				"tag": "cats with images",
				"id": "1273026480692322304"
			}],
			"meta": {
				"sent": "2020-06-16T22:55:39.356Z",
				"summary": {
					"created": 1,
					"not_created": 0,
					"valid": 1,
					"invalid": 0
				}
			}
		}`))
		require.NoError(t, err)
	})

	response, _, err := service.AddRule(context.Background(), AddCommand{
		Add: []Add{{
			Value: "cat has:images",
			Tag:   "cats with images",
		}},
	}, false)

	require.NoError(t, err)
	require.EqualValues(t, 1, response.Meta.Summary.Created)
}

func TestGetRules(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/2/tweets/search/stream/rules", req.URL.Path)
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [{
				"id": "1273028376882589696",
				"value": "cat has:images",
				"tag": "cats with images"
			}],
			"meta": {
				"sent": "2020-06-16T23:14:06.498Z"
			}
		}`))
		require.NoError(t, err)
	})

	response, _, err := service.GetRules(context.Background())

	require.NoError(t, err)
	require.Equal(t, "cat has:images", response.Data[0].Value)
	require.EqualValues(t, "2020-06-16T23:14:06.498Z", response.Meta.Sent.Format("2006-01-02T15:04:05.000Z"))
}
