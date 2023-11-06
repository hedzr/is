package term

import (
	"io"
	"os"
	"runtime"
)

var (
	EnvironmentOverrideColors bool
	DisableColors             bool // no colorful paint? false to disable colorful paint.
	ForceColors               bool // always colorful paint? true to enable colorful paint even if the underlying tty cannot support ANSI escaped sequences.
)

// IsTty detects a writer if it is abstracting from a tty (console, terminal) device.
func IsTty(w io.Writer) bool {
	var isTerminal = checkIfTerminal(w)
	return isTerminal
}

// IsColored detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func IsColored(w io.Writer) bool {
	isColored := ForceColors || (IsTty(w) && (runtime.GOOS != "windows"))

	if EnvironmentOverrideColors {
		if force, ok := os.LookupEnv("FORCE_COLOR"); ok && force != "0" {
			isColored = true
		} else if ok && force == "0" {
			isColored = false
		} else if os.Getenv("COLOR") == "0" {
			isColored = false
		}
	}

	return isColored && !DisableColors
}
