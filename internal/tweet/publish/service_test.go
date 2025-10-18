package publish

import (
    "context"
    "encoding/json"
    "io"
    "net/http"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestCreateTweet(t *testing.T) {
    service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
        require.Equal(t, http.MethodPost, req.Method)
        require.Equal(t, "/2/tweets", req.URL.Path)

        body, err := io.ReadAll(req.Body)
        require.NoError(t, err)
        require.JSONEq(t, `{"text":"hello"}`, string(body))

        res.Header().Set("x-rate-limit-limit", "300")
        res.Header().Set("x-rate-limit-remaining", "299")
        res.Header().Set("x-rate-limit-reset", "42")
        res.WriteHeader(http.StatusCreated)
        err = json.NewEncoder(res).Encode(CreateTweetResponse{Data: TweetData{ID: "1", Text: "hello"}})
        require.NoError(t, err)
    })

    response, rateLimits, err := service.CreateTweet(context.Background(), CreateTweetRequest{Text: "hello"})
    require.NoError(t, err)
    require.Equal(t, "1", response.Data.ID)
    require.Equal(t, 300, rateLimits.Limit)
    require.Equal(t, 299, rateLimits.Remaining)
    require.Equal(t, 42, rateLimits.Reset)
}

func TestDeleteTweet(t *testing.T) {
    service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
        require.Equal(t, http.MethodDelete, req.Method)
        require.Equal(t, "/2/tweets/1", req.URL.Path)
        res.WriteHeader(http.StatusOK)
        err := json.NewEncoder(res).Encode(DeleteTweetResponse{Data: DeleteData{Deleted: true}})
        require.NoError(t, err)
    })

    response, _, err := service.DeleteTweet(context.Background(), "1")
    require.NoError(t, err)
    require.True(t, response.Data.Deleted)
}
