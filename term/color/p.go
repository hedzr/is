package color

import (
	"fmt"
	"os"
	"strings"

	"github.com/hedzr/is/states"
)

//nolint:lll //no
func _internalLogTo(tofn func(sb strings.Builder, ln bool), format string, args ...any) { //nolint:goprintffuncname //so what
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(format, args...))
	tofn(sb, strings.HasSuffix(sb.String(), "\n"))
}

// // Error outputs formatted message to stderr.
// func Error(format string, args ...any) { //nolint:goprintffuncname //so what
// 	_internalLogTo(func(sb strings.Builder, ln bool) {
// 		if ln {
// 			_, _ = fmt.Fprint(os.Stderr, sb.String())
// 		} else {
// 			// _, _ = fmt.Fprintln(os.Stderr, sb.String())
// 			lx.Errorf("%v", sb.String())
// 		}
// 	}, format, args...)
// }
//
// // Fatal outputs formatted message to stderr.
// func Fatal(format string, args ...any) { //nolint:goprintffuncname //so what
// 	_internalLogTo(func(sb strings.Builder, ln bool) {
// 		if ln {
// 			lx.Fatalf("%v", sb.String())
// 		} else {
// 			lx.Panicf("%v", sb.String())
// 		}
// 	}, format, args...)
// }
//
// // Warn outputs formatted message to stderr while logger level
// // less than level.WarnLevel.
// // For level.SetLevel(level.ErrorLevel), the text will be discarded.
// func Warn(format string, args ...any) { //nolint:goprintffuncname //so what
// 	// for the key scene who want quiet output, we may disable
// 	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
// 	if level.GetLevel() < level.WarnLevel {
// 		return
// 	}
//
// 	_internalLogTo(func(sb strings.Builder, ln bool) {
// 		if ln {
// 			print(ToColor(FgYellow, sb.String())) //nolint:forbidigo //no
// 		} else {
// 			// println(sb.String())
// 			lx.Warnf("%v", ToColor(FgYellow, sb.String()))
// 		}
// 	}, format, args...)
// }
//
// // Log will print the formatted message to stdout.
// //
// // While the message ends with '\n', it will be printed by
// // print(), so you'll see it always.
// // But if not, the message will be printed by hedzr/log. In
// // this case, its outputting depends on
// // hedzr/level.GetLogLevel() >= level.DebugLevel.
// //
// // Log outputs formatted message to stdout while logger level
// // less than level.WarnLevel.
// // For level.SetLevel(level.ErrorLevel), the text will be discarded.
// func Log(format string, args ...any) { //nolint:goprintffuncname //so what
// 	// for the key scene who want quiet output, we may disable
// 	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
// 	if level.GetLevel() < level.WarnLevel {
// 		return
// 	}
//
// 	_internalLogTo(func(sb strings.Builder, ln bool) {
// 		if ln {
// 			print(sb.String()) //nolint:forbidigo //no
// 		} else {
// 			// println(sb.String())
// 			lx.Printf("%v", sb.String())
// 		}
// 	}, format, args...)
// }
//
// // Verbose outputs formatted message to stdout while cmdr is in
// // VERBOSE mode.
// // For level.SetLevel(level.ErrorLevel), the text will be discarded.
// func Verbose(format string, args ...any) { //nolint:goprintffuncname //so what
// 	if states.Env().IsVerboseMode() {
// 		_internalLogTo(func(sb strings.Builder, ln bool) {
// 			if ln {
// 				print(sb.String()) //nolint:forbidigo //no
// 			} else {
// 				// println(sb.String())
// 				lx.Printf("%v", sb.String())
// 			}
// 		}, format, args...)
// 	}
// }
//
// // Trace outputs formatted message to stdout while logger level
// // is level.TraceLevel, or cmdr is in TRACE mode or trace module
// // is enabled.
// func Trace(format string, args ...any) { //nolint:goprintffuncname //so what
// 	if level.GetLevel() == level.TraceLevel || !states.Env().GetTraceMode() {
// 		return
// 	}
//
// 	_internalLogTo(func(sb strings.Builder, ln bool) {
// 		if ln {
// 			// log.Skip(extrasLogSkip).Tracef("%v", sb.String())
// 			Colored(FgLightGray, "%v", sb.String())
// 		} else {
// 			lx.Tracef("%v", sb.String())
// 		}
// 	}, format, args...)
// }

// Highlight outputs formatted message to stdout while logger level
// less than level.WarnLevel.
// For level.SetLevel(level.ErrorLevel), the text will be discarded.
func Highlight(format string, args ...any) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
	// if level.GetLevel() < level.WarnLevel {
	// 	return
	// }

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if states.Env().IsNoColorMode() {
			_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "\x1b[0;1m%v\x1b[0m", sb.String())
		}
		if !ln {
			println() //nolint:forbidigo //no
		}
	}, format, args...)
}

// Dimf outputs formatted message to stdout while logger level less
// than level.WarnLevel and cmdr is in verbose mode.
//
// For example, after level.SetLevel(level.ErrorLevel), the text via Dimf will be discarded.
//
// While env-var VERBOSE=1, the text via Dimf will be shown.
func Dimf(format string, args ...any) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
	// if level.GetLevel() < level.WarnLevel {
	// 	return
	// }

	if states.Env().IsVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			if states.Env().IsNoColorMode() {
				_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
			} else {
				_, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())
			}
			if !ln {
				println() //nolint:forbidigo //no
			}
		}, format, args...)
	}
}

// Text prints formatted message without any predefined ansi escaping.
func Text(format string, args ...any) { //nolint:goprintffuncname //so what
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

// Dim outputs formatted message to stdout while logger level
// less than level.WarnLevel.
//
// For example, after level.SetLevel(level.ErrorLevel), the text via Dim will be discarded.
func Dim(format string, args ...any) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
	// if level.GetLevel() < level.WarnLevel {
	// 	return
	// }

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if states.Env().IsNoColorMode() {
			_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
		} else {
			_, _ = fmt.Fprintf(os.Stdout, "\x1b[2m\x1b[37m%v\x1b[0m", sb.String())
		}
		if !ln {
			println() //nolint:forbidigo //no
		}
	}, format, args...)
}

func ToDim(format string, args ...any) (str string) {
	str = fmt.Sprintf(format, args...)
	if states.Env().IsNoColorMode() {
		return
	}
	str = fmt.Sprintf("\x1b[2m\x1b[37m%v\x1b[0m", str)
	return
}

func ToHighlight(format string, args ...any) (str string) {
	str = fmt.Sprintf(format, args...)
	if states.Env().IsNoColorMode() {
		return
	}
	str = fmt.Sprintf("\x1b[0;1m%v\x1b[0m", str)
	return
}

func ToColor(clr Color, format string, args ...any) (str string) {
	// str = fmt.Sprintf(format, args...)
	// if states.Env().IsNoColorMode() {
	// 	return
	// }
	// str = fmt.Sprintf("\u001B[%dm%v\x1b[0m", clr, str)

	var sb strings.Builder
	text := fmt.Sprintf(format, args...)
	WrapColorTo(&sb, clr, text)
	return sb.String()
}

// Coloredf outputs formatted message to stdout while logger level
// less than level.WarnLevel and cmdr is in VERBOSE mode.
func Coloredf(clr Color, format string, args ...any) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
	// if level.GetLevel() < level.WarnLevel {
	// 	return
	// }

	if states.Env().IsVerboseMode() {
		_internalLogTo(func(sb strings.Builder, ln bool) {
			if states.Env().IsNoColorMode() {
				_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
			} else {
				color(clr)
				_, _ = fmt.Fprintf(os.Stdout, "%v\x1b[0m", sb.String())
			}
			if !ln {
				println() //nolint:forbidigo //no
			}
		}, format, args...)
	}
}

// Colored outputs formatted message to stdout while logger level
// less than level.WarnLevel.
// For level.SetLevel(level.ErrorLevel), the text will be discarded.
func Colored(clr Color, format string, args ...any) { //nolint:goprintffuncname //so what
	// for the key scene who want quiet output, we may disable
	// most of the messages by cmdr.SetLogLevel(level.ErrorLevel)
	// if level.GetLevel() < level.WarnLevel {
	// 	return
	// }

	_internalLogTo(func(sb strings.Builder, ln bool) {
		if states.Env().IsNoColorMode() {
			_, _ = fmt.Fprintf(os.Stdout, "%v", sb.String())
		} else {
			color(clr)
			_, _ = fmt.Fprintf(os.Stdout, "%v\x1b[0m", sb.String())
		}
		if !ln {
			println() //nolint:forbidigo //no
		}
	}, format, args...)
}

func color(c Color) {
	_, _ = fmt.Fprintf(os.Stdout, "\x1b[%dm", c)
}

func ResetColor(c Color) { //nolint:unused //no
	_, _ = fmt.Fprint(os.Stdout, "\x1b[0m")
}
