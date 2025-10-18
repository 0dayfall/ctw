package user

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlock(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/2/users/1/blocking", req.URL.Path)
		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"target_user_id":"2"}`, string(body))
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(`{"data": {"blocking": true}}`))
		require.NoError(t, err)
	})

	resp, _, err := service.Block(context.Background(), "1", "2")
	require.NoError(t, err)
	require.True(t, resp.Data.Blocking)
}

func TestUnblock(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodDelete, req.Method)
		require.Equal(t, "/2/users/1/blocking/2", req.URL.Path)
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"data": {"blocking": false}}`))
		require.NoError(t, err)
	})

	resp, _, err := service.Unblock(context.Background(), "1", "2")
	require.NoError(t, err)
	require.False(t, resp.Data.Blocking)
}

func TestFollow(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/2/users/1/following", req.URL.Path)
		body, err := io.ReadAll(req.Body)
		require.NoError(t, err)
		require.JSONEq(t, `{"target_user_id":"2"}`, string(body))
		res.WriteHeader(http.StatusOK)
		_, err = res.Write([]byte(`{"data": {"following": true, "pending_follow": false}}`))
		require.NoError(t, err)
	})

	resp, _, err := service.Follow(context.Background(), "1", "2")
	require.NoError(t, err)
	require.True(t, resp.Data.Following)
	require.False(t, resp.Data.PendingFollow)
}

func TestUnfollow(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodDelete, req.Method)
		require.Equal(t, "/2/users/1/following/2", req.URL.Path)
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"data": {"following": false}}`))
		require.NoError(t, err)
	})

	resp, _, err := service.Unfollow(context.Background(), "1", "2")
	require.NoError(t, err)
	require.False(t, resp.Data.Following)
}
