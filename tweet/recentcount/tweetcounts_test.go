package tweetcounts

import (
	"os"
	"testing"

	"github.com/0dayfall/carboncopy/httphandler"
)

func TestMain(m *testing.M) {
	httphandler.Init("AAAAAAAAAAAAAAAAAAAAAI0kRgEAAAAAypS5hDlUu0fQxhPfsegcVRKgGyE%3Dz5LtZJLTBtN5xgrCCGAvQPX8a8fZFxkKJVhWCHkHkIEoSDCcvM")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCountRecent(t *testing.T) {
	tweerCount := GetRecentCount("FB lang:en", "day")
	httphandler.PrettyPrint(tweerCount)
}
