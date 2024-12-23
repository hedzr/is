package is

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"golang.org/x/term"

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

func TestTty(t *testing.T) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is being piped to stdin")
	} else {
		fmt.Println("stdin is from a terminal")
	}

	stat, _ = os.Stdout.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is being piped to Stdout")

		for i, k := range os.Environ() {
			if l := strings.ToLower(k); strings.Contains(l, "term") || strings.Contains(l, "color") {
				t.Logf("%5d. %s", i+1, k)
			}
		}
	} else {
		fmt.Println("Stdout is from a terminal")
	}

	t.Run("testIsTerminal", testIsTerminal)
}

// testIsTerminal cannot work properly because os.Stdout was been
// piped to testing framework, that is, it is never a terminal/tty
// device.
func testIsTerminal(t *testing.T) {
	v(t, "os.Stdout", os.Stdout)
	v(t, "os.Stderr", os.Stderr)

	// is terminal
	run4dev(t, "/dev/ptmx", v)
	run4dev(t, "/dev/pty", v)
	run4dev(t, "/dev/tty", v)

	// is tty
	isTty(t, "os.Stdout", os.Stdout)
	run4dev(t, "/dev/ptmx", isTty)

	// is colorful
	isColorful(t, "os.Stdout", os.Stdout)
	run4dev(t, "/dev/ptmx", isColorful)

	// window size
	winSize(t, "os.Stdout", os.Stdout)
	run4dev(t, "/dev/tty", winSize)

	// _, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	// _, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), ioctlReadTermios)
	// _, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TIOCGETA)
	// if err != nil {
	// 	fmt.Println("Hello World")
	// } else {
	// 	fmt.Println("\033[1;34mHello World!\033[0m")
	// }

	t.Log("------------")

	if runtime.GOOS == "darwin" {
		f := os.Stdout
		// t.Logf("stdout is terminal: %v", terminal.IsTerminal(int(f.Fd())))

		stat, _ := f.Stat()
		mod := stat.Mode()
		t.Logf("mod: %v (pipe: %v)", mod.Type(), mod&os.ModeNamedPipe != 0)
		if mod&os.ModeNamedPipe != 0 {
			// if !terminal.IsTerminal(0) || !terminal.IsTerminal(1) {
			// 	t.Errorf("stdin/stdout should be terminal")
			// 	t.FailNow()
			// }
			//
			// term.NewTerminal(c, "Fxx: ")
		}
		return
	}
}

func run4dev(t *testing.T, what string, cb func(t *testing.T, what string, f *os.File)) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skipf("unknown terminal path for GOOS %v", runtime.GOOS)
	}

	file, err := os.OpenFile(what, os.O_RDWR, 0)
	if err != nil {
		t.Logf("%q cannot be opened, skip. (err: %v)", what, err)
		return
	}
	defer file.Close()

	cb(t, what, file)
}

func isTty(t *testing.T, what string, f *os.File) {
	t.Logf("IsTty(%q): %v (mine)",
		what, Tty(f),
	)
}

func isColorful(t *testing.T, what string, f *os.File) {
	t.Logf("IsColored(%q): %v (mine)",
		what, ColoredTty(f),
	)
}

func winSize(t *testing.T, what string, f *os.File) {
	c, r := GetTtySize()
	t.Logf("Window Size (%q): %v cols x %v rows",
		what, c, r,
	)
}

func v(t *testing.T, what string, f *os.File) {
	fd := int(f.Fd())
	t.Logf("IsTerminal(%q): %v / %v (mine), %v (ignore), %v (golang.org/x/term)",
		what, Terminal(f), TerminalFd(f.Fd()),
		"", // terminal.IsTerminal(fd),
		term.IsTerminal(fd),
	)
}

func TestTtyFunctions(t *testing.T) {
	t.Logf(`TTY: %v, colorful: %v
	string is escaped(plain): %v
	string is escaped(color): %v
	string is escaped(color): %v
	IsStartupByDoubleClick:   %v
	`,
		Tty(os.Stdout), ColoredTty(os.Stdout),
		TtyEscaped("plain text"),
		TtyEscaped("\x1b[2mcolor\x1b0m text"),
		AnsiEscaped("\x1b[2mcolor\x1b0m text"),
		StartupByDoubleClick(),
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
	x, y, _ = GetTtySizeByName("/dev/stdout")
	t.Logf("%v, %v", x, y)
	x, y, _ = GetTtySizeByFile(os.Stdout)
	t.Logf("%v, %v", x, y)
	x, y, _ = GetTtySizeByFd(os.Stdout.Fd())
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

func TestToBool(t *testing.T) {
	assertEqual(t, true, ToBool(1))
	assertEqual(t, true, ToBool("yes"))
	assertEqual(t, true, ToBool("TRUE"))
	assertEqual(t, true, StringToBool("on"))
	assertEqual(t, true, StringToBool("y"))
	assertEqual(t, true, StringToBool("t"))
	assertEqual(t, true, StringToBool("YES"))
	assertEqual(t, true, StringToBool("true"))
}

//

//

func assertEqual(t testing.TB, expect, actual interface{}, messages ...any) {
	assertEqualSkip(t, 2, expect, actual, messages...)
}

func assertEqualSkip(t testing.TB, skip int, expect, actual interface{}, messages ...any) {
	if !isEqual(expect, actual) {
		_, file, line, _ := runtime.Caller(skip)
		if len(messages) > 0 {
			fmt.Println(messages...)
		}
		fmt.Printf("%s:%d expecting %v (%v), but got %v (%v).\n",
			path.Base(file), line,
			expect, reflect.TypeOf(expect),
			actual, reflect.TypeOf(actual),
			// diffmatchpatch.DiffValues(expect, actual),
		)
		t.FailNow()
	}
}

func isEqual(val1, val2 any) bool {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if !v1.IsValid() && !v2.IsValid() {
		return true
	}

	if v1.Kind() == reflect.Slice && v2.Kind() == reflect.Slice {
		return isEqualArray(v1, v2)
	}

	// return reflect.DeepEqual(val1, val2)
	return risEqualConcretely(v1, v2)
}

func isEqualArray(v1, v2 reflect.Value) bool {
	if v1.Len() != v2.Len() {
		return false
	}

	for i := 0; i < v1.Len(); i++ {
		x1 := v1.Index(i)
		var found bool
		for j := 0; j < v2.Len(); j++ {
			x2 := v2.Index(j)
			if risEqualConcretely(x1, x2) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// isEqualConcretely returns whether val1 is equal to val2 taking into account Pointers, Interfaces and their underlying types
func isEqualConcretely(val1, val2 any) bool {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)

	if v1.Kind() == reflect.Ptr {
		v1 = v1.Elem()
	}

	if v2.Kind() == reflect.Ptr {
		v2 = v2.Elem()
	}

	if !v1.IsValid() && !v2.IsValid() {
		return true
	}

	return risEqualConcretely(v1, v2)
}

func risEqualConcretely(v1, v2 reflect.Value) bool {
	switch v1.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if v1.IsNil() {
			v1 = reflect.ValueOf(nil)
		}
	}

	switch v2.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		if v2.IsNil() {
			v2 = reflect.ValueOf(nil)
		}
	}

	v1Underlying := reflect.Zero(reflect.TypeOf(v1)).Interface()
	v2Underlying := reflect.Zero(reflect.TypeOf(v2)).Interface()

	if v1 == v1Underlying {
		if v2 == v2Underlying {
			goto CASE4
		} else {
			goto CASE3
		}
	} else {
		if v2 == v2Underlying {
			goto CASE2
		} else {
			goto CASE1
		}
	}

CASE1:
	return reflect.DeepEqual(v1.Interface(), v2.Interface())
CASE2:
	return reflect.DeepEqual(v1.Interface(), v2)
CASE3:
	return reflect.DeepEqual(v1, v2.Interface())
CASE4:
	return reflect.DeepEqual(v1, v2)
}

func TestDetectors(t *testing.T) {
	t.Logf("os.args: %v", os.Args)
	assertEqual(t, true, InTesting())
	t.Logf(`
	InDebugging:         %v, DebuggerAttached: %v, InDebugMode: %v,
	InTracing:           %v,
	InTestingT(os.Args): %v,
	InDevelopingTime:    %v,
	InDockerEnvSimple:   %v,
	InVscodeTerminal:    %v,
	InK8s:               %v / %v,
	InIstio:             %v,
	InDocker:            %v,
	DockerBuild:         %v,
	K8sBuild:            %v,
	IstioBuild:          %v,
	DebugBuild:          %v,
	VerboseBuild:        %v,
	VerboseModeEnabled:  %v,
	GetVerboseLevel:     %v,
	QuietModeEnabled:    %v,
	GetQuietLevel:       %v,
	NoColorMode:         %v,
	GetNoColorLevel:     %v,
	DebugMode:           %v,
	GetDebugLevel:       %v,
	TraceMode:           %v,
	GetTraceLevel:       %v,
	
	States(): %v,
	Env():    %v,
`,
		InDebugging(), DebuggerAttached(), InDebugMode(),
		InTracing(),
		InTestingT(os.Args),
		InDevelopingTime(),
		InDockerEnvSimple(),
		InVscodeTerminal(),

		InK8s(), InK8sYN(),
		InIstio(),
		InDocker(),
		DockerBuild(),
		K8sBuild(),
		IstioBuild(),

		DebugBuild(),
		VerboseBuild(),

		VerboseModeEnabled(), GetVerboseLevel(),
		QuietModeEnabled(), GetQuietLevel(),
		NoColorMode(), GetNoColorLevel(),
		DebugMode(), GetDebugLevel(),
		TraceMode(), GetTraceLevel(),

		States(),
		Env(),
	)
}
