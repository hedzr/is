package is

import (
	"io"
	"os"

	term2 "golang.org/x/term"

	"github.com/hedzr/is/term"
	"github.com/hedzr/is/term/chk"
	"github.com/hedzr/is/term/color"
)

// Tty detects a writer if it is abstracting from a tty (console, terminal) device.
func Tty(w io.Writer) bool { return chk.IsTty(w) }

// ColoredTty detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func ColoredTty(w io.Writer) bool { return chk.IsColored(w) }

// ColorfulTty detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func ColorfulTty(w io.Writer) bool { return chk.IsColorful(w) }

// TtyEscaped detects a string if it contains ansi color escaped sequences
// Deprecated v0.5.3, use HasAnsiEscaped
func TtyEscaped(s string) bool { return term.IsAnsiEscaped(s) }

// AnsiEscaped detects a string if it contains ansi color escaped sequences
func AnsiEscaped(s string) bool { return term.IsAnsiEscaped(s) }

// StripEscapes removes any ansi color escaped sequences from a string
func StripEscapes(str string) (strCleaned string) { return term.StripEscapes(str) }

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) {
	var b []byte
	b, err = term.ReadPassword(int(os.Stdin.Fd()))
	text = string(b)
	return
}

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) { return term.GetTtySize() }

// GetTtySizeByName retrieve terminal window size by device name. such as "/dev/tty"
func GetTtySizeByName(fn string) (cols, rows int, err error) { return term.GetTtySizeByName(fn) }

// GetTtySizeByFile retrieve terminal window size by *os.File object. such as os.Stdout
func GetTtySizeByFile(f *os.File) (cols, rows int, err error) { return term.GetTtySizeByFile(f) }

// GetTtySizeByFd retrieve terminal window size by fd (file-descriptor). such as [os.Stdout.Fd()]
func GetTtySizeByFd(fd uintptr) (cols, rows int, err error) { return term.GetTtySizeByFd(fd) }

// StartupByDoubleClick detects
// if windows golang executable file is running via double click or from cmd/shell terminator
func StartupByDoubleClick() bool { return term.IsStartupByDoubleClick() }

// CygwinTerminal return true if the file descriptor is a cygwin or msys2
// terminal. This is also always false on this environment.
func CygwinTerminal(fd uintptr) bool {
	return false
}

// Terminal returns whether the given file descriptor is a terminal.
// Terminal detects if a file is a terminal device (tty)
func Terminal(f *os.File) bool {
	ret := term2.IsTerminal(int(f.Fd()))
	return ret
}

// TerminalFd detects if a file-descriptor is a terminal device (tty)
func TerminalFd(fd uintptr) bool {
	ret := term.IsTerminal(int(fd))
	return ret
}

// Color returns an indexer for term/color subpackage.
//
// For example, call the Translator to convert the html tags to color codes in a string:
//
//	is.Color().GetColorTranslator().Translate("<b>bold</b>")
func Color() color.Index { return color.Index{} }

func StdoutStat() (normalFile, redirected, piped, terminal bool) {
	return chk.StatStdout()
}

func StdoutPiped() bool {
	_, _, b, _ := chk.StatStdout()
	return b
}

func StdoutRedirected() bool {
	_, b, _, _ := chk.StatStdout()
	return b
}

func StdoutIsTerminal() bool {
	_, _, _, b := chk.StatStdout()
	return b
}

func StdoutIsNormalFile() bool {
	b, _, _, _ := chk.StatStdout()
	return b
}
