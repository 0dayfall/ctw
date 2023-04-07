package tweet

import (
	"os"
	"testing"

	"github.com/0dayfall/ctw/internal/utils"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCountRecent(t *testing.T) {
	tweerCount := GetRecentCount("FB lang:en", "day")
	utils.PrettyPrint(tweerCount)
}
