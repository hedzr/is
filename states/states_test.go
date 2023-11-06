package states

import (
	"testing"
)

func TestMinimalEnv_CountOfQuiet(t *testing.T) {
	t.Logf("%v", Env().CountOfQuiet())
	Env().SetQuietCount(8)
	t.Logf("%v", Env().CountOfQuiet())
}
