package is

import (
	"runtime/debug"
	"testing"

	"github.com/hedzr/is/stringtool"
)

func TestRandomStringPure(t *testing.T) {
	t.Log(stringtool.RandomStringPure(8))
}

func TestDebugBuildInfo(t *testing.T) {
	if info, ok := debug.ReadBuildInfo(); ok {
		t.Logf("info: %+v", info)
	}
}
