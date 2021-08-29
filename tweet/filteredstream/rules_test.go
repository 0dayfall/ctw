package tweet

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
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
