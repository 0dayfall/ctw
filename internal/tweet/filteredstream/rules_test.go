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
	rulesUrl = testServer.URL
	response, err := AddRule(AddCommand{
		Add: []Add{{
			Value: "cat has:images",
			Tag:   "cats with images",
		}},
	}, true)
	require.NoError(t, err)
	require.Equal(t, "cat has:images", response.Data[0].Value)
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
	rulesUrl = testServer.URL
	response, err := AddRule(AddCommand{
		Add: []Add{{
			Value: "cat has:images",
			Tag:   "cats with images",
		}},
	}, false)
	require.NoError(t, err)
	require.Equal(t, "cat has:images", response.Data[0].Value)
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
	rulesUrl = testServer.URL
	response, err := GetRules()
	require.NoError(t, err)
	require.Equal(t, "cat has:images", response.Data[0].Value)
}

func makeMap() map[string]string {
	m := make(map[string]string)
	m["tweet.fields"] = "created_at"
	m["expansions"] = "author_id"
	m["user.fields"] = "created_at"
	return m
}

func TestStreamsURL(t *testing.T) {
	require.Equal(t, "https://api.twitter.com/2/tweets/search/stream?"+
		"expansions=author_id&tweet.fields=created_at&user.fields=created_at",
		createStreamUrlWithFields(makeMap()))
}

func TestStream(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{
			"data": [
			  {
				"author_id": "2244994945",
				"created_at": "2022-09-14T19:00:55.000Z",
				"id": "1228393702244134912",
				"edit_history_tweet_ids": ["1228393702244134912"],
				"text": "What did the developer write in their Valentine’s card?\n  \nwhile(true) {\n    I = Love(You);  \n}"
			  },
			  {
				"author_id": "2244994945",
				"created_at": "2022-09-12T17:09:56.000Z",
				"id": "1227640996038684673",
				 "edit_history_tweet_ids": ["1227640996038684673"],
				"text": "Doctors: Googling stuff online does not make you a doctor\n\nDevelopers: https://t.co/mrju5ypPkb"
			  },
			  {
				"author_id": "2244994945",
				"created_at": "2022-09-27T20:26:41.000Z",
				"id": "1199786642791452673",
				"edit_history_tweet_ids": ["1199786642791452673"],
				"text": "C#"
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
	}))
	defer func() { testServer.Close() }()
	rulesUrl = testServer.URL
	response, err := Stream(makeMap())
	require.NoError(t, err)
	// Check that the name field of the first user is "Twitter Dev"
	if includes, ok := response.(map[string]interface{})["includes"].(map[string]interface{}); ok {
		if users, ok := includes["users"].([]interface{}); ok && len(users) > 0 {
			if user, ok := users[0].(map[string]interface{}); ok {
				if name, ok := user["name"].(string); ok {
					if name != "Twitter Dev" {
						t.Errorf("Expected name to be \"Twitter Dev\", but got %q", name)
					}
				} else {
					t.Errorf("Expected name to be a string, but got %T", user["name"])
				}
			} else {
				t.Error("Expected first user to be a map")
			}
		} else {
			t.Error("Expected includes.users to be a non-empty array")
		}
	} else {
		t.Error("Expected includes to be a map")
	}
}
