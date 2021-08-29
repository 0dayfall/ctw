package tweet

import (
	"os"
	"testing"

	"github.com/0dayfall/carboncopy/config"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSearchRecent(t *testing.T) {
	config.Init("AAAAAAAAAAAAAAAAAAAAAI0kRgEAAAAAypS5hDlUu0fQxhPfsegcVRKgGyE%3Dz5LtZJLTBtN5xgrCCGAvQPX8a8fZFxkKJVhWCHkHkIEoSDCcvM")
	SearchRecent("ericsson lang:sv")
}

func TestSearchRecentNextToken(t *testing.T) {
	config.Init("AAAAAAAAAAAAAAAAAAAAAI0kRgEAAAAAypS5hDlUu0fQxhPfsegcVRKgGyE%3Dz5LtZJLTBtN5xgrCCGAvQPX8a8fZFxkKJVhWCHkHkIEoSDCcvM")
	_, _, token := SearchRecent("ericsson lang:sv")
	_, _, _ = SearchRecentNextToken("ericsson lang:sv", token)
}
