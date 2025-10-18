package tweet

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStreamSendsQueryParameters(t *testing.T) {
	fields := map[string]string{
		"tweet.fields": "created_at",
		"expansions":   "author_id",
		"user.fields":  "created_at",
	}

	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/2/tweets/search/stream", req.URL.Path)
		for key, value := range fields {
			require.Equal(t, value, req.URL.Query().Get(key))
		}

		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [
				{
					"author_id": "2244994945",
					"created_at": "2022-09-14T19:00:55.000Z",
					"id": "1228393702244134912",
					"text": "What did the developer write in their Valentine’s card?"
				}
			],
			"includes": {
				"users": [
					{
						"created_at": "2013-12-14T04:35:55.000Z",
						"id": "2244994945",
						"name": "Twitter Dev",
						"username": "TwitterDev"
					}
				]
			}
		}`))
		require.NoError(t, err)
	})

	response, _, err := service.Stream(context.Background(), fields)
	require.NoError(t, err)
	require.Len(t, response.Data, 1)
	require.Equal(t, "1228393702244134912", response.Data[0].ID)
	require.Len(t, response.Includes.Users, 1)
	require.Equal(t, "Twitter Dev", response.Includes.Users[0].Name)
}
