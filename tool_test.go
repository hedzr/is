package is

import (
	"os"
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

func TestTtyFunctions(t *testing.T) {
	t.Logf(`TTY: %v, colorful: %v
	string is escaped(plain): %v
	string is escaped(color): %v
	`,
		Tty(os.Stdout), ColoredTty(os.Stdout),
		TtyEscaped("plain text"),
		TtyEscaped("\x1b[2mcolor\x1b0m text"),
	)

	t.Logf("%v", StripEscapes(`
	
		<code>code</code>
	NC Cool
	 But it's tight.
	  Hold On!
	Hurry Up.
	`))

	x, y := GetTtySize()
	t.Logf("%v, %v", x, y)
}
