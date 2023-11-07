package is

import (
	"io"

	"github.com/hedzr/is/term"
)

// Tty detects a writer if it is abstracting from a tty (console, terminal) device.
func Tty(w io.Writer) bool { return term.IsTty(w) }

// ColoredTty detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func ColoredTty(w io.Writer) bool { return term.IsColored(w) }

// TtyEscaped detects a string if it contains ansi color escaped sequences
func TtyEscaped(s string) bool { return term.IsTtyEscaped(s) }

// StripEscapes removes any ansi color escaped sequences from a string
func StripEscapes(str string) (strCleaned string) { return term.StripEscapes(str) }

// ReadPassword reads the password from stdin with safe protection
func ReadPassword() (text string, err error) { return term.ReadPassword() }

// GetTtySize returns the window size in columns and rows in the active console window.
// The return value of this function is in the order of cols, rows.
func GetTtySize() (cols, rows int) { return term.GetTtySize() }
