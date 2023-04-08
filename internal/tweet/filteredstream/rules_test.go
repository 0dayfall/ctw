package tweet

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestAddRuleDryRun(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
	}))
	defer func() { testServer.Close() }()
	rulesBaseUrl = testServer.URL
	response, err := AddRule(AddCommand{
		Add: []Add{{
			Value: "cat has:images",
			Tag:   "cats with images",
		}},
	}, true)
	require.NoError(t, err)
	require.EqualValues(t, 1, response.Meta.Summary.Created)
}

func TestAddRule(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
	}))
	defer func() { testServer.Close() }()
	rulesBaseUrl = testServer.URL
	response, err := AddRule(AddCommand{
		Add: []Add{{
			Value: "cat has:images",
			Tag:   "cats with images",
		}},
	}, false)
	require.NoError(t, err)
	require.EqualValues(t, 1, response.Meta.Summary.Created)
}

func TestGetRules(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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
	}))
	defer func() { testServer.Close() }()
	rulesBaseUrl = testServer.URL
	response, err := GetRules()
	require.NoError(t, err)
	require.Equal(t, "cat has:images", response.Data[0].Value)
	require.EqualValues(t, "2020-06-16T23:14:06.498Z", response.Meta.Sent.Format("2006-01-02T15:04:05.000Z"))
}
