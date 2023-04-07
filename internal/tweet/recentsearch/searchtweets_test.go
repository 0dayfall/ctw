package tweet

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSearchRecent(t *testing.T) {
	SearchRecent("ericsson lang:sv")
}

func TestSearchRecentNextToken(t *testing.T) {
	_, _, token := SearchRecent("ericsson lang:sv")
	_, _, _ = SearchRecentNextToken("ericsson lang:sv", token)
}
