package term

import (
	"io"

	"github.com/hedzr/is/term/chk"
)

// IsTty detects a writer if it is abstracting from a tty (console, terminal) device.
func IsTty(w io.Writer) bool { return chk.IsTty(w) }

// IsColored detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func IsColored(w io.Writer) bool { return chk.IsColorful(w) }

// IsColorful detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func IsColorful(w io.Writer) (colorful bool) { return chk.IsColorful(w) }

func StatStdout() (normalFile, redirected, piped, term bool) { return chk.StatStdout() }

func StatStdoutString() (status string) { return chk.StatStdoutString() }
func StdoutIsPiped() (b bool)           { return chk.StdoutIsPiped() }
