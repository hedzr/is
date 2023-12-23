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

func TestFileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "1*.log")
	if err != nil {
		t.Error(err)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write([]byte("OK\r\n"))
	if err != nil {
		t.Error(err)
	}
	tmpFile.Close()

	if !FileExists(tmpFile.Name()) {
		t.Fatalf("expecting temp file is existed: %q", tmpFile.Name())
	}
	if c, err := ReadFile(tmpFile.Name()); err != nil {
		t.Error(err)
	} else if string(c) != "OK\r\n" {
		t.Fatalf("file content unmatched. read: %v", c)
	}
}
