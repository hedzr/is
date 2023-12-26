package buildtags

import (
	"testing"
)

func TestIsBuildTagExists(t *testing.T) {
	b := IsBuildTagExists("verbose")
	t.Logf("verbose: %v", b)
}
