package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLookupID(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/2/users/123", req.URL.Path)
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"data": {"id": "123", "name": "name", "username": "user"}}`))
		require.NoError(t, err)
	})

	user, _, err := service.LookupID(context.Background(), "123", nil)
	require.NoError(t, err)
	require.Equal(t, "123", user.ID)
}

func TestLookupUsername(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/2/users/by/username/jane", req.URL.Path)
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"data": {"id": "1", "name": "Jane", "username": "jane"}}`))
		require.NoError(t, err)
	})

	user, _, err := service.LookupUsername(context.Background(), "jane", nil)
	require.NoError(t, err)
	require.Equal(t, "jane", user.UserName)
}

func TestLookupIDs(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/2/users", req.URL.Path)
		require.Equal(t, "1,2", req.URL.Query().Get("ids"))
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"data": [{"id": "1", "username": "a"}, {"id": "2", "username": "b"}]}`))
		require.NoError(t, err)
	})

	users, _, err := service.LookupIDs(context.Background(), []string{"1", "2"}, nil)
	require.NoError(t, err)
	require.Len(t, users, 2)
}

func TestLookupUsernames(t *testing.T) {
	service := newTestService(t, func(res http.ResponseWriter, req *http.Request) {
		require.Equal(t, "/2/users/by", req.URL.Path)
		require.Equal(t, "alice,bob", req.URL.Query().Get("usernames"))
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"data": [{"id": "1", "username": "alice"}]}`))
		require.NoError(t, err)
	})

	users, _, err := service.LookupUsernames(context.Background(), []string{"alice", "bob"}, nil)
	require.NoError(t, err)
	require.Len(t, users, 1)
}
