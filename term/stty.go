package term

import (
	"io"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var (
	EnvironmentOverrideColors bool  = true
	DisableColors             bool  // no colorful paint? false to disable colorful paint.
	ForceColors               bool  // always colorful paint? true to enable colorful paint even if the underlying tty cannot support ANSI escaped sequences.
	MinVal                    int64 // 0,16,88,256, or 1<<24
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
func IsColored(w io.Writer) bool { return IsColorful(w) }

// IsColorful detects a writer if it is a colorful tty device.
//
// A colorful tty device can receive ANSI escaped sequences and draw its.
func IsColorful(w io.Writer) (colorful bool) {
	// && (runtime.GOOS != "windows")

	// Check for Azure DevOps pipelines.
	// Has to be above the `!streamIsTTY` check.
	if _, colorful = os.LookupEnv("TF_BUILD"); colorful {
		return
	}
	if _, colorful = os.LookupEnv("AGENT_NAME"); colorful {
		return
	}

	// var term string
	// var ok bool
	var force, disabled bool

	colorful, force, disabled = IsTty(w), ForceColors, DisableColors

	if runtime.GOOS == "windows" {
		// Windows 10 build 10586 is the first Windows release that supports 256 colors.
		// Windows 10 build 14931 is the first release that supports 16m/TrueColor.
		MinVal = 1 << 24
		force = checkForceColor()
		disabled = checkDisableColor()
		return !DisableColors && (MinVal > 0 || force)
	}

	if EnvironmentOverrideColors {
		force = checkForceColor()
		disabled = checkDisableColor()
		_ = disabled

		var name, val string
		var present bool

		if name, val, present = anyInEnv("CI", "CI_RUNNING"); present {
			if name, val, present = anyInEnv(
				"GITHUB_ACTIONS", "GITEA_ACTIONS", "CIRCLECI",
			); present {
				MinVal = 1 << 24
			}
			if name, val, present = anyInEnv(
				"TRAVIS", "APPVEYOR", "GITLAB_CI", "BUILDKITE", "DRONE",
			); present {
				MinVal = 16
			}
		}
		_ = name

		if name, val, present = anyInEnv("COLORTERM"); present {
			if val == "truecolor" {
				MinVal = 1 << 24
			} else {
				MinVal = 16
			}
		}
		if name, val, present = anyInEnv("TERM"); present {
			switch val {
			case "dumb":
				MinVal = 16
			default:
				if strings.HasPrefix(val, "xterm-") {
					MinVal = 256
				} else if val == "xterm-kitty" {
					MinVal = 1 << 24
				} else if regexp.MustCompile(`^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux`).Match([]byte(val)) {
					MinVal = 16
				}
			}
			MinVal = 1 << 24
		}
	}

	// return force || isColorful && !disabled
	return !DisableColors && (MinVal > 0 || force)
}

func checkDisableColor() (disabled bool) {
	if name, val, present := anyInEnv("NO_COLOR", "NOCOLOR"); present {
		DisableColors = StringToBool(val)
		disabled = DisableColors
		_ = name
	}
	return
}

func checkForceColor() (force bool) {
	if b, ok := os.LookupEnv("FORCE_COLOR"); ok {
		if v, e := strconv.ParseInt(b, 10, 64); e == nil {
			minval := v
			switch minval {
			case 1:
				minval = 16
			case 2:
				minval = 88
			case 3:
				minval = 256
			case 4:
				minval = 1 << 24
			}
			MinVal = minval
		}
		force = MinVal > 0
	}
	return
}

func anyInEnv(names ...string) (name, v string, yes bool) {
	for _, name = range names {
		if v, yes = os.LookupEnv(name); yes {
			return
		}
	}
	return
}

}
