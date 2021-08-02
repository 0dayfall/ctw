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
