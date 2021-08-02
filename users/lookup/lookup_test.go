package lookup

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

func TestLookupUsernames(t *testing.T) {
	LookupUsernames([]string{"greenorangebay", "eleandreon"})
}

func TestLookupUsername(t *testing.T) {
	LookupUsername("greenorangebay")
}

func TestLookupID(t *testing.T) {
	LookupID("greenorangebay")
}

func TestLookupIDs(t *testing.T) {
	LookupIDs([]string{"1232142", "1232434"})
}
