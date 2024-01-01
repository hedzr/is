package term

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

func TestIsTerminal(t *testing.T) {
	v(t, "os.Stdout", os.Stdout)
	v(t, "os.Stderr", os.Stderr)

	// is terminal
	vdev(t, "/dev/ptmx")
	vdev(t, "/dev/pty")
	vdev(t, "/dev/tty")

	// is tty
	isTty(t, "os.Stdout", os.Stdout)
	tdev(t, "/dev/ptmx")

	// is colorful
	isColorful(t, "os.Stdout", os.Stdout)
	cdev(t, "/dev/ptmx")

	// window size
	winSize(t, "os.Stdout", os.Stdout)
	run4dev(t, "/dev/tty", winSize)

	// _, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	// _, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), ioctlReadTermios)
	_, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TIOCGETA)
	if err != nil {
		fmt.Println("Hello World")
	} else {
		fmt.Println("\033[1;34mHello World!\033[0m")
	}

	cols, rows, err := getDeviceSize(os.Stdout)
	t.Logf("getDeviceSize: %v, %v, %v", cols, rows, err)
}

func TestIsTerminalTempFile2(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "TestIsTerminalTempFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if IsTerminal(int(file.Fd())) {
		t.Fatalf("IsTerminal unexpectedly returned true for temporary file %s", file.Name())
	}
}

func TestIsTerminalTerm2(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		t.Skipf("unknown terminal path for GOOS %v", runtime.GOOS)
	}
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		// t.Fatal(err)
		t.Logf("[WARN] /dev/ptmx under darwin failed. err: %v", err)
		return
	}
	defer file.Close()

	if !term.IsTerminal(int(file.Fd())) {
		// t.Fatalf("IsTerminal unexpectedly returned false for terminal file %s", file.Name())
		t.Logf("[WARN] /dev/ptmx under darwin failed. err: %v", err)
	}
}

func tdev(t *testing.T, what string) { run4dev(t, what, isTty) }
func vdev(t *testing.T, what string) { run4dev(t, what, v) }
func cdev(t *testing.T, what string) { run4dev(t, what, isColorful) }

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
		what, IsTty(f),
	)
}

func isColorful(t *testing.T, what string, f *os.File) {
	t.Logf("IsColored(%q): %v (mine)",
		what, IsColored(f),
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
	t.Logf("IsTerminal(%q): %v (mine), %v (golang.org/x/crypto/ssh/terminal), %v (golang.org/x/term)",
		what, IsTerminal(fd),
		terminal.IsTerminal(fd), term.IsTerminal(fd),
	)
}
