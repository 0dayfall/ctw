package tweet

import (
	"os"
	"testing"

	"github.com/0dayfall/carboncopy/config"
)

func TestMain(m *testing.M) {
	config.Init("AAAAAAAAAAAAAAAAAAAAAI0kRgEAAAAAypS5hDlUu0fQxhPfsegcVRKgGyE%3Dz5LtZJLTBtN5xgrCCGAvQPX8a8fZFxkKJVhWCHkHkIEoSDCcvM")
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestAddRuleDryRun(t *testing.T) {
	AddRule(AddCommand{
		Add: []Add{{
			Value: "Test value",
			Tag:   "Test tag",
		}},
	}, true)
}

func TestGetRules(t *testing.T) {
	GetRules()
}

func TestStream(t *testing.T) {
	Stream()
}
