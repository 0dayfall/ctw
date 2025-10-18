package user

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0dayfall/ctw/internal/client"
	"github.com/stretchr/testify/require"
)

func newTestService(t *testing.T, handler http.HandlerFunc) *Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg := client.Config{
		BaseURL:     server.URL + "/",
		BearerToken: "test-token",
	}

	c, err := client.New(cfg)
	require.NoError(t, err)

	return NewService(c)
}
