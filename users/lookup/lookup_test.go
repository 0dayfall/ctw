package lookup

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
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
