package color

import (
	"testing"

	"github.com/hedzr/is/term"
)

func TestIsAnsiEscaped(t *testing.T) {
	for i, tc := range []struct {
		src    string
		expect bool
	}{
		{"", false},
		{GetCPT().Translate(`<code>code</code>`, FgDefault), true},
	} {
		got := term.IsAnsiEscaped(tc.src)
		if got != tc.expect {
			t.Fatalf("%5d. IsAnsiEscaped(%q) failed: expecting '%v' but got '%v'", i, tc.src, tc.expect, got)
		}
	}
}

func TestHighlight(t *testing.T) {
	Highlight("Highlight: hello, %v!", "world")
}

func TestDimV(t *testing.T) {
	Dimf("Dimf (verbose build only): hello, %v!", "world")
}

func TestText(t *testing.T) {
	Text("Text: hello, %v!\n", "world")
}

func TestDim(t *testing.T) {
	Dim("Dim: hello, %v!", "world")
}

func TestToDim(t *testing.T) {
	t.Logf("%v", ToDim("ToDim: hello, %v!", "world"))
}

func TestToColor(t *testing.T) {
	t.Logf("%v", ToColor(FgMagenta, "ToColor: hello, %v!", "world"))
}

func TestColoredV(t *testing.T) {
	Coloredf(FgLightMagenta, "Coloredf (verbose build only): hello, %v!", "world")
}

func TestColored(t *testing.T) {
	Colored(FgLightMagenta, "Colored: hello, %v!", "world")
}
