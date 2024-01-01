package term

import (
	"io"
	"os"

	"golang.org/x/term"
)

var (
	EnvironmentOverrideColors bool = true
	DisableColors             bool // no colorful paint? false to disable colorful paint.
	ForceColors               bool // always colorful paint? true to enable colorful paint even if the underlying tty cannot support ANSI escaped sequences.
)

// IsTty detects a writer if it is abstracting from a tty (console, terminal) device.
func IsTty(w io.Writer) bool {
	switch z := w.(type) {
	case *os.File:
		return term.IsTerminal(int(z.Fd()))
	default:
		return false
	}
}

// IsColored detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func IsColored(w io.Writer) bool {
	// && (runtime.GOOS != "windows")

	isColored, force := IsTty(w), ForceColors

	if EnvironmentOverrideColors {
		if b, ok := os.LookupEnv("FORCE_COLOR"); ok {
			force = StringToBool(b)
		}

		if b, ok := os.LookupEnv("NO_COLOR"); ok && StringToBool(b) {
			DisableColors = true
		} else if b, ok := os.LookupEnv("NOCOLOR"); ok && StringToBool(b) {
			DisableColors = true
		}
	}

	return force || isColored && !DisableColors
}
